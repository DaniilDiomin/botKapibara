package handler

import (
	"bot4/products"
	states "bot4/state"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"sort"
)

// Структура для хранения состояния опроса
type surveyState struct {
	AllProducts []struct {
		Category string
		Product  products.ProductItem
	}
	CurrentIndex int
}

func (h *Handler) handleRoleSelection(chatID int64, msg *tgbotapi.Message) {
	restaraunt, err := h.stateManager.GetContext(chatID, "pending_restaurant")
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка: не выбран ресторан"))
		h.stateManager.SetState(chatID, states.Idle)
		return
	}

	switch msg.Text {
	case "Повар":
		switch restaraunt {
		case "rogachev":
			h.startSurvey(chatID, h.products.Kapibara.Cook, restaraunt, "cook")
		case "rechica":
			h.startSurvey(chatID, h.products.Kapibara.Cook, restaraunt, "cook")
		}
	case "Кассир":
		switch restaraunt {
		case "rogachev":
			h.startSurvey(chatID, h.products.Kapibara.Cashier, restaraunt, "cashier")
		case "rechica":
			h.startSurvey(chatID, h.products.Kapibara.Cashier, restaraunt, "cashier")
		}
	default:
		h.bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, выберите роль из предложенных"))
	}
}

func (h *Handler) startSurvey(chatID int64, productsMap map[string][]products.ProductItem, restaurant, role string) {
	log.Printf("[startSurvey] Начало опроса для chatID: %d, ресторан: %s, роль: %s", chatID, restaurant, role)

	// Проверяем наличие данных о товарах
	if productsMap == nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Нет данных о товарах"))
		return
	}

	// Собираем только непустые категории
	categories := make([]string, 0)
	for cat, items := range productsMap {
		if len(items) > 0 {
			categories = append(categories, cat)
		}
	}
	sort.Strings(categories)

	if len(categories) == 0 {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Нет доступных товаров"))
		return
	}

	// Формируем плоский список всех товаров
	allProducts := h.getAllProducts(categories, productsMap)
	log.Printf("[startSurvey] Всего товаров: %d", len(allProducts))

	// Создаем начальное состояние опроса
	state := surveyState{
		AllProducts:  allProducts,
		CurrentIndex: 0,
	}

	// Сохраняем состояние
	h.stateManager.SetContext(chatID, "survey_state", h.serializeSurveyState(state))
	h.stateManager.SetContext(chatID, "survey_restaurant", restaurant)
	h.stateManager.SetContext(chatID, "survey_role", role)
	h.stateManager.SetState(chatID, states.SurveyInProgress)
	log.Printf("[startSurvey] Состояние сохранено")

	// Отправляем первый вопрос
	h.sendNextQuestion(chatID, state)
}

func (h *Handler) sendNextQuestion(chatID int64, state surveyState) {
	log.Printf("[sendNextQuestion] Отправка вопроса для chatID: %d, индекс: %d/%d",
		chatID, state.CurrentIndex, len(state.AllProducts))

	// Проверяем, не закончились ли товары
	if state.CurrentIndex >= len(state.AllProducts) {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Опрос завершен!"))
		h.stateManager.SetState(chatID, states.Idle)
		h.stateManager.SetContext(chatID, "survey_state", "")
		log.Printf("[sendNextQuestion] Опрос завершен для chatID: %d", chatID)
		return
	}

	// Получаем текущий товар
	current := state.AllProducts[state.CurrentIndex]
	log.Printf("[sendNextQuestion] Текущий товар: %s - %s", current.Category, current.Product.Name)

	// Создаем временный update для клавиатуры
	fakeUpdate := &tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	// Формируем текст вопроса
	questionText := fmt.Sprintf("%s - %s\nВыберите количество:", current.Category, current.Product.Name)

	// Отправляем клавиатуру
	if err := h.createKeyboard(h.bot, fakeUpdate, current.Product.Quantity, questionText); err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка создания клавиатуры: "+err.Error()))
		log.Printf("[sendNextQuestion] Ошибка создания клавиатуры: %v", err)
	}
}

func (h *Handler) createKeyboard(bot *tgbotapi.BotAPI, update *tgbotapi.Update, options []string, text string) error {
	// Проверяем, что слайс не пустой
	if len(options) == 0 {
		return fmt.Errorf("options slice is empty")
	}

	options = append(options, "Пропустить")

	// Создаем клавиатуру
	var keyboardRows [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton

	// Формируем кнопки (по 2 кнопки в ряд)
	for i, option := range options {
		row = append(row, tgbotapi.NewKeyboardButton(option))
		// Если набрали 2 кнопки или это последняя опция
		if len(row) == 2 || i == len(options)-1 {
			keyboardRows = append(keyboardRows, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}

	// Создаем объект клавиатуры
	replyKeyboard := tgbotapi.NewReplyKeyboard(keyboardRows...)
	replyKeyboard.OneTimeKeyboard = true // Автоматическое скрытие после выбора

	// Настраиваем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = replyKeyboard

	// Отправляем сообщение с клавиатурой
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func (h *Handler) handleSurveyResponse(chatID int64, msg *tgbotapi.Message) {
	log.Printf("[handleSurveyResponse] Получен ответ от chatID: %d, текст: %s", chatID, msg.Text)

	// Получаем текущее состояние опроса
	state, err := h.getSurveyState(chatID)
	if err != nil {
		log.Printf("[handleSurveyResponse] Ошибка получения состояния: %v", err)
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка: "+err.Error()))
		h.stateManager.SetState(chatID, states.Idle)
		return
	}
	log.Printf("[handleSurveyResponse] Текущий индекс: %d/%d", state.CurrentIndex, len(state.AllProducts))

	// Проверяем, не закончились ли товары
	if state.CurrentIndex >= len(state.AllProducts) {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Опрос завершен!"))
		h.stateManager.SetState(chatID, states.Idle)
		h.stateManager.SetContext(chatID, "survey_state", "")
		log.Printf("[handleSurveyResponse] Опрос завершен для chatID: %d", chatID)
		return
	}

	// Получаем текущий товар
	current := state.AllProducts[state.CurrentIndex]
	log.Printf("[handleSurveyResponse] Обработка товара: %s - %s", current.Category, current.Product.Name)

	// Обрабатываем ответ (если не "Пропустить")
	if msg.Text != "Пропустить" {
		// Получаем ресторан и роль для форматирования сообщения
		restaurant, _ := h.stateManager.GetContext(chatID, "survey_restaurant")
		role, _ := h.stateManager.GetContext(chatID, "survey_role")

		// Форматируем сообщение в зависимости от ресторана
		messageText := fmt.Sprintf("%s: %s %s", restaurant, current.Category, current.Product.Name)
		if role != "" {
			messageText = fmt.Sprintf("%s - %s", current.Category, current.Product.Name)
		}
		messageText += " " + msg.Text

		log.Printf("[handleSurveyResponse] Отправка в группу: %s", messageText)
		if err := h.sendTextInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WriteoffTopicID, messageText); err != nil {
			log.Printf("[handleSurveyResponse] Ошибка отправки в группу: %v", err)
		}
	}

	// Увеличиваем индекс
	newState := surveyState{
		AllProducts:  state.AllProducts,
		CurrentIndex: state.CurrentIndex + 1,
	}
	log.Printf("[handleSurveyResponse] Новый индекс: %d", newState.CurrentIndex)

	// Сохраняем новое состояние
	h.stateManager.SetContext(chatID, "survey_state", h.serializeSurveyState(newState))

	// Отправляем следующий вопрос
	h.sendNextQuestion(chatID, newState)
}

func (h *Handler) getAllProducts(categories []string, productsMap map[string][]products.ProductItem) []struct {
	Category string
	Product  products.ProductItem
} {
	var allProducts []struct {
		Category string
		Product  products.ProductItem
	}

	for _, cat := range categories {
		productsInCat := productsMap[cat]
		for _, prod := range productsInCat {
			allProducts = append(allProducts, struct {
				Category string
				Product  products.ProductItem
			}{cat, prod})
		}
	}
	return allProducts
}

// Сериализует состояние опроса в строку
func (h *Handler) serializeSurveyState(state surveyState) string {
	bytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("Ошибка сериализации состояния: %v", err)
		return ""
	}
	return string(bytes)
}

// Десериализует состояние опроса из строки
func (h *Handler) deserializeSurveyState(data string) (surveyState, error) {
	var state surveyState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return surveyState{}, fmt.Errorf("ошибка десериализации: %v", err)
	}
	return state, nil
}

// Получает состояние опроса из StateManager
func (h *Handler) getSurveyState(chatID int64) (surveyState, error) {
	stateStr, err := h.stateManager.GetContext(chatID, "survey_state")
	if err != nil {
		return surveyState{}, fmt.Errorf("состояние опроса не найдено: %v", err)
	}
	return h.deserializeSurveyState(stateStr)
}

// Функция отправки текста в топик
func (h *Handler) sendTextInTopic(chatID, themeID int64, text string) error {
	// Формируем параметры запроса
	params := map[string]interface{}{
		"chat_id":              chatID,
		"text":                 text,
		"message_thread_id":    themeID,
		"disable_notification": false,
	}

	// Сериализуем параметры в JSON
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("ошибка сериализации параметров: %v", err)
	}

	// Создаем HTTP-запрос
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.config.Token)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем ответ
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка Telegram API (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
