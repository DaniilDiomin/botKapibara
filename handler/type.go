package handler

import (
	"bot4/config"
	"bot4/products"
	"bot4/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot          *tgbotapi.BotAPI
	stateManager *state.StateManager
	products     *products.ProductsConfig
	config       *config.Config
}
