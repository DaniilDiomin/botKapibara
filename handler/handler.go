package handler

import (
	"bot4/config"
	"bot4/products"
	"bot4/state"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	case state.RequestSubmission:
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
		//
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
		//
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
		//
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
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toRechicaSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toRechicaSelection(chatID, message)
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
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toRogachevSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Kapibara.GroupChatID, h.config.Kapibara.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toRogachevSelection(chatID, message)
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
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "1")
		case "Окончил смену":
			forwardMsg := fmt.Sprintf("@%s %s", message.Chat.UserName, message.Text)
			if err := h.sendMsgInTopic(h.config.Freshkof.GroupChatID, h.config.Freshkof.WorkHoursTopicID, forwardMsg); err != nil {
				log.Println(err)
			}
			h.toFreshcoffSelection(chatID, message)
			h.stateManager.SetContext(chatID, "is_work", "")
		case "Назад":
			h.toFreshcoffSelection(chatID, message)
		}
	}
}

func (h *Handler) sendMsgInTopic(chatID, themeID int64, msg string) error {
	type sendMessageRequest struct {
		ChatID          int64  `json:"chat_id"`
		MessageThreadID int64  `json:"message_thread_id"`
		Text            string `json:"text"`
	}

	reqBody := sendMessageRequest{
		ChatID:          chatID,
		MessageThreadID: themeID,
		Text:            msg,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.bot.Token)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %d", resp.StatusCode)
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
		case "Продукт":
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
		case "Продукт":
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
		case "Продукт":
		case "Назад":
			h.toFreshcoffSelection(chatID, message)
		}

	}

}
