package handler

import (
	"bot4/config"
	"bot4/products"
	"bot4/state"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHandler(bot *tgbotapi.BotAPI, stateManager *state.StateManager, products *products.ProductsConfig, config *config.Config) *Handler {
	return &Handler{
		bot:          bot,
		stateManager: stateManager,
		products:     products,
		config:       config,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	userState := h.stateManager.GetState(chatID)

	if update.Message.IsCommand() {
		h.handleCommand(chatID, update.Message)
		return
	}

	if update.Message.Photo != nil {
		if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WriteoffTopicID, update.Message); err != nil {
			log.Println(err)
		}
	}

	switch userState.Current {
	case state.Idle:
		h.handleIdle(chatID)
	case state.RestarauntSelection:
		h.handleRestarauntSelection(chatID, update.Message)
	case state.FreshcoffSelection:
		h.handleFreshcoffSelection(chatID, update.Message)
	case state.RechicaSelection:
		h.handleRechicaSelection(chatID, update.Message)
	case state.RogachevSelection:
		h.handleRogachevSelection(chatID, update.Message)
	case state.WorkSchedule:
		h.handleWorkSchedule(chatID, update.Message)
	case state.WriteOff:
		h.handleWriteOff(chatID, update.Message)
	case state.RequestSubmission: h.handleRequestSubmission(chatID, update.Message)
	}
}

func (h *Handler) handleCommand(chatID int64, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		h.stateManager.SetState(chatID, state.Idle)
		msg := tgbotapi.NewMessage(chatID, "Привет! Выберите место:")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("КАПИБАРА РОГАЧЕВ"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("КАПИБАРА РЕЧИЦА"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("ФРЕШКОФФ РОГАЧЕВ"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)
		h.stateManager.SetState(chatID, state.RestarauntSelection)
	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
		h.bot.Send(msg)
	}
}

func (h *Handler) handleIdle(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Нажмите /start, чтобы начать.")
	h.bot.Send(msg)
}

func (h *Handler) handleRestarauntSelection(chatID int64, message *tgbotapi.Message) {
	switch message.Text {
	case "КАПИБАРА РОГАЧЕВ":
		h.stateManager.SetContext(chatID, "restaraunt", "rogachev")
		h.stateManager.SetState(chatID, state.RogachevSelection)

		msg := tgbotapi.NewMessage(chatID, "Вы выбрали КАПИБАРА РОГАЧЕВ. Что дальше?")
		keyboard := tgbotapi.NewReplyKeyboard(
			// Каждая кнопка в отдельном ряду
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Учет рабочего времени"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заявка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Списание"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)
	case "КАПИБАРА РЕЧИЦА":
		h.stateManager.SetContext(chatID, "restaraunt", "rechica")
		h.stateManager.SetState(chatID, state.RechicaSelection)
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали КАПИБАРА РЕЧИЦА. Что дальше?")

		keyboard := tgbotapi.NewReplyKeyboard(
			// Каждая кнопка в отдельном ряду
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Учет рабочего времени"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заявка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Списание"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)
	case "ФРЕШКОФФ РОГАЧЕВ":
		h.stateManager.SetContext(chatID, "restaraunt", "freshcoff")
		h.stateManager.SetState(chatID, state.FreshcoffSelection)
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали ФРЕШКОФФ РОГАЧЕВ. Что дальше?")

		keyboard := tgbotapi.NewReplyKeyboard(
			// Каждая кнопка в отдельном ряду
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Учет рабочего времени"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заявка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Списание"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)
	}
}

func (h *Handler) toRestarauntSelection(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetState(chatID, state.RestarauntSelection)
	h.stateManager.SetContext(chatID, "restaraunt", "")
	msg := tgbotapi.NewMessage(chatID, "Привет! Выберите место:")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("КАПИБАРА РОГАЧЕВ"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("КАПИБАРА РЕЧИЦА"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ФРЕШКОФФ РОГАЧЕВ"),
		),
	)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) handleRogachevSelection(chatID int64, message *tgbotapi.Message) {
	switch message.Text {
	case "Учет рабочего времени":
		h.handleWorkSchedule(chatID, message)
	case "Списание":
		h.handleWriteOff(chatID, message)
	case "Заявка":
		h.handleRequestSubmission(chatID, message)
	case "Назад":
		h.toRestarauntSelection(chatID, message)
	default:
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите опцию из предложенных.")
		h.bot.Send(msg)
	}
}
func (h *Handler) toRogachevSelection(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetContext(chatID, "restaraunt", "rogachev")
	h.stateManager.SetState(chatID, state.RogachevSelection)

	msg := tgbotapi.NewMessage(chatID, "готово.")
	keyboard := tgbotapi.NewReplyKeyboard(
		// Каждая кнопка в отдельном ряду
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Учет рабочего времени"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заявка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Списание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) handleFreshcoffSelection(chatID int64, message *tgbotapi.Message) {
	switch message.Text {
	case "Учет рабочего времени":
		h.handleWorkSchedule(chatID, message)
	case "Списание":
		h.handleWriteOff(chatID, message)
	case "Заявка":
		h.handleRequestSubmission(chatID, message)
	case "Назад":
		h.toRestarauntSelection(chatID, message)
	default:
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите опцию из предложенных.")
		h.bot.Send(msg)
	}
}
func (h *Handler) toFreshcoffSelection(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetContext(chatID, "restaraunt", "rogachev")
	h.stateManager.SetState(chatID, state.FreshcoffSelection)

	msg := tgbotapi.NewMessage(chatID, "готово.")
	keyboard := tgbotapi.NewReplyKeyboard(
		// Каждая кнопка в отдельном ряду
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Учет рабочего времени"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заявка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Списание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) handleRechicaSelection(chatID int64, message *tgbotapi.Message) {
	switch message.Text {
	case "Учет рабочего времени":
		h.handleWorkSchedule(chatID, message)
	case "Списание":
		h.handleWriteOff(chatID, message)
	case "Заявка":
		h.handleRequestSubmission(chatID, message)
	case "Назад":
		h.toRestarauntSelection(chatID, message)
	default:
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите опцию из предложенных.")
		h.bot.Send(msg)
	}
}
func (h *Handler) toRechicaSelection(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetContext(chatID, "restaraunt", "rogachev")
	h.stateManager.SetState(chatID, state.RechicaSelection)

	msg := tgbotapi.NewMessage(chatID, "готово.")
	keyboard := tgbotapi.NewReplyKeyboard(
		// Каждая кнопка в отдельном ряду
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Учет рабочего времени"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заявка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Списание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) handleWorkSchedule(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetState(chatID, state.WorkSchedule)
	restaraunt, err := h.stateManager.GetContext(chatID, "restaraunt")
	if err != nil {
		log.Println(err)
	}
	switch restaraunt {
	case "rechica":
		msg := tgbotapi.NewMessage(chatID, "Статус работы:")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("На работе"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Окончил смену"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "На работе":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "")
		case "Назад":
			h.toRechicaSelection(chatID, message)
		}

	case "rogachev":
		msg := tgbotapi.NewMessage(chatID, "Статус работы:")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("На работе"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Окончил смену"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "На работе":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "")
		case "Назад":
			h.toRogachevSelection(chatID, message)
		}

	case "freshcoff":
		msg := tgbotapi.NewMessage(chatID, "Статус работы:")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("На работе"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Окончил смену"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "На работе":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, message); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "")
		case "Назад":
			h.toFreshcoffSelection(chatID, message)
		}
	}
}

func (h *Handler) forwardMsgInTopic(chatID, themeID int64, message *tgbotapi.Message) error {
	// Проверяем, есть ли последнее сообщение пользователя
	if message == nil || message.Chat == nil || message.Chat.ID == 0 {
		return fmt.Errorf("нет доступного сообщения для пересылки")
	}

	// Формируем параметры запроса
	params := map[string]interface{}{
		"chat_id":              chatID,
		"from_chat_id":         message.Chat.ID,
		"message_id":           message.MessageID,
		"message_thread_id":    themeID,
		"disable_notification": false,
	}

	// Сериализуем параметры в JSON
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("ошибка сериализации параметров: %v", err)
	}

	// Создаем HTTP-запрос
	url := fmt.Sprintf("https://api.telegram.org/bot%s/forwardMessage", h.config.Token)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем ответ
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка Telegram API (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (h *Handler) handleWriteOff(chatID int64, message *tgbotapi.Message) {
	h.stateManager.SetState(chatID, state.WriteOff)
	restaraunt, err := h.stateManager.GetContext(chatID, "restaraunt")
	if err != nil {
		log.Println(err)
	}
	switch restaraunt {
	case "rogachev":
		msg := tgbotapi.NewMessage(chatID, "Что списать?")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заготовка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Продукт"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "Заготовка":
			fallthrough
		case "Продукт":
			h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Что списываем? Отправьте фото и объясните причину"))
			if message.Text != "Заготовка" && message.Text != "Продукт" {
				if err := h.forwardMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WriteoffTopicID, message); err != nil {
					log.Println(err)
				}
			}
		case "Назад":
			h.toRogachevSelection(chatID, message)
		}

	case "rechica":
		msg := tgbotapi.NewMessage(chatID, "Что списать?")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заготовка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Продукт"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "Заготовка":
			fallthrough
		case "Продукт":
			h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Что списываем? Отправьте фото и объясните причину"))
			if message.Text != "Заготовка" && message.Text != "Продукт" {
				if err := h.forwardMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WriteoffTopicID, message); err != nil {
					log.Println(err)
				}
			}
		case "Назад":
			h.toRechicaSelection(chatID, message)
		}

	case "freshcoff":
		msg := tgbotapi.NewMessage(chatID, "Что списать?")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Заготовка"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Продукт"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		)
		msg.ReplyMarkup = keyboard
		h.bot.Send(msg)

		switch message.Text {
		case "Заготовка":
			fallthrough
		case "Продукт":
			h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Что списываем? Отправьте фото и объясните причину"))
			if message.Text != "Заготовка" && message.Text != "Продукт" {
				if err := h.forwardMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WriteoffTopicID, message); err != nil {
					log.Println(err)
				}
			}
		case "Назад":
			h.toFreshcoffSelection(chatID, message)
		}

	}
}

func (h *Handler) handleRequestSubmission(chatID int64, msg *tgbotapi.Message) {
	h.stateManager.SetState(chatID, state.RequestSubmission)
	restaraunt, err := h.stateManager.GetContext(chatID, "restaraunt")
	if err != nil {
		log.Println(err)
	}
	switch restaraunt {
	case "rogachev":
	case "rechica":
	case "freshcoff":
	}
}
