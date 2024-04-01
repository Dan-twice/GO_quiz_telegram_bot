package init

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
)

const (
	botToken = "botToken"
	baseURL  = "https://heroku_app_name_registered.herokuapp.com/"
)

var bot *tgbotapi.BotAPI

func Telegram() {
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	// bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := baseURL + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Println(err)
	}
}
