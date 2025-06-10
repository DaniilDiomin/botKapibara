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
	"strconv"
	"strings"
)

// Структура для хранения состояния опроса
type surveyState struct {
	Categories   []string
	CurrentIndex int // Единый индекс для всех товаров
}

func (h *Handler) requestSubmissionFreshcoff(chatID int64, msg *tgbotapi.Message) {
	log.Printf("[requestSubmissionFreshcoff] Начало опроса для chatID: %d", chatID)

	// Проверяем наличие данных о товарах
	if h.products == nil || h.products.Fresfcoff == nil || h.products.Fresfcoff.Products == nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Нет данных о товарах"))
		log.Printf("[requestSubmissionFreshcoff] Ошибка: нет данных о товарах")
		return
	}

	// Собираем только непустые категории
	categories := make([]string, 0)
	for cat, items := range h.products.Fresfcoff.Products {
		if len(items) > 0 {
			categories = append(categories, cat)
		}
	}
	sort.Strings(categories)

	if len(categories) == 0 {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Нет доступных товаров"))
		log.Printf("[requestSubmissionFreshcoff] Ошибка: нет доступных товаров")
		return
	}

	// Создаем начальное состояние опроса
	state := surveyState{
		Categories:   categories,
		CurrentIndex: 0,
	}

	// Сохраняем состояние
	h.stateManager.SetContext(chatID, "survey_state", h.serializeSurveyState(state))
	h.stateManager.SetState(chatID, states.SurveyInProgress)
	log.Printf("[requestSubmissionFreshcoff] Состояние сохранено: %s", h.serializeSurveyState(state))

	// Отправляем первый вопрос
	h.sendNextQuestion(chatID, state)
}

func (h *Handler) sendNextQuestion(chatID int64, state surveyState) {
	log.Printf("[sendNextQuestion] Отправка вопроса для chatID: %d, индекс: %d", chatID, state.CurrentIndex)

	// Получаем все товары в плоском списке
	allProducts := h.getAllProducts(state.Categories)
	log.Printf("[sendNextQuestion] Всего товаров: %d", len(allProducts))

	// Проверяем, не закончились ли товары
	if state.CurrentIndex >= len(allProducts) {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Опрос завершен!"))
		h.stateManager.SetState(chatID, states.Idle)
		h.stateManager.SetContext(chatID, "survey_state", "")
		log.Printf("[sendNextQuestion] Опрос завершен для chatID: %d", chatID)
		return
	}

	// Получаем текущий товар
	current := allProducts[state.CurrentIndex]
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
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text+":")
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
	log.Printf("[handleSurveyResponse] Текущее состояние: %s", h.serializeSurveyState(state))

	// Получаем все товары в плоском списке
	allProducts := h.getAllProducts(state.Categories)
	log.Printf("[handleSurveyResponse] Всего товаров: %d", len(allProducts))

	// Проверяем, не закончились ли товары
	if state.CurrentIndex >= len(allProducts) {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Опрос завершен!"))
		h.stateManager.SetState(chatID, states.FreshcoffSelection)
		h.stateManager.SetContext(chatID, "survey_state", "")
		log.Printf("[handleSurveyResponse] Опрос завершен для chatID: %d", chatID)
		h.toFreshcoffSelection(chatID, msg)
		return
	}

	// Получаем текущий товар
	current := allProducts[state.CurrentIndex]
	log.Printf("[handleSurveyResponse] Обработка товара: %s - %s", current.Category, current.Product.Name)

	// Обрабатываем ответ (если не "Пропустить")
	if msg.Text != "Пропустить" {
		messageText := fmt.Sprintf("%s - %s", current.Product.Name, msg.Text)
		log.Printf("[handleSurveyResponse] Отправка в группу: %s", messageText)
		if err := h.sendTextInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.ProcurementTopicID, messageText); err != nil {
			log.Printf("[handleSurveyResponse] Ошибка отправки в группу: %v", err)
		}
	}

	// Увеличиваем индекс
	newState := surveyState{
		Categories:   state.Categories,
		CurrentIndex: state.CurrentIndex + 1,
	}
	log.Printf("[handleSurveyResponse] Новое состояние: %s", h.serializeSurveyState(newState))

	// Сохраняем новое состояние
	h.stateManager.SetContext(chatID, "survey_state", h.serializeSurveyState(newState))

	// Отправляем следующий вопрос
	h.sendNextQuestion(chatID, newState)
}

func (h *Handler) getAllProducts(categories []string) []struct {
	Category string
	Product  products.ProductItem
} {
	var allProducts []struct {
		Category string
		Product  products.ProductItem
	}

	for _, cat := range categories {
		productsInCat := h.products.Fresfcoff.Products[cat]
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
	return fmt.Sprintf("%s|%d", strings.Join(state.Categories, ","), state.CurrentIndex)
}

// Десериализует состояние опроса из строки
func (h *Handler) deserializeSurveyState(data string) (surveyState, error) {
	parts := strings.Split(data, "|")
	if len(parts) < 2 {
		return surveyState{}, fmt.Errorf("неверный формат состояния: %s", data)
	}

	categories := strings.Split(parts[0], ",")
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return surveyState{}, fmt.Errorf("ошибка преобразования индекса: %v", err)
	}

	return surveyState{
		Categories:   categories,
		CurrentIndex: index,
	}, nil
}

// Получает состояние опроса из StateManager
func (h *Handler) getSurveyState(chatID int64) (surveyState, error) {
	stateStr, err := h.stateManager.GetContext(chatID, "survey_state")
	if err != nil {
		return surveyState{}, fmt.Errorf("состояние опроса не найдено: %v", err)
	}

	state, err := h.deserializeSurveyState(stateStr)
	if err != nil {
		return surveyState{}, fmt.Errorf("ошибка десериализации: %v", err)
	}

	return state, nil
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
