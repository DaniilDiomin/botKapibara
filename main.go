package main

import (
	"bot4/config"
	"bot4/handler"
	"bot4/products"
	"bot4/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	DEBUG = true
)

func main() {
	//Получаем конфиг
	cfg, err := config.LoadConfig("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	//Получаем продукты
	productCfg, err := products.LoadProducts("productsKapibara.json", "productsFreshcoff.json")
	if err != nil {
		log.Fatal(err)
	}

	//создаем бота
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = DEBUG

	stateManager := state.NewStateManager()

	h := handler.NewHandler(bot, stateManager, productCfg, cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//канад со всеми обновлениями
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		h.HandleUpdate(update)
	}

	//fmt.Println(productCfg.Fresfcoff.Products["КОФЕ"][0].Name)
}
