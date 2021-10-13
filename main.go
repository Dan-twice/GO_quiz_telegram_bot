package main

import (
	_ "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	rd "math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// type ContextKey string
// const ContextUserKey ContextKey = "user"

var bot *tgbotapi.BotAPI
var userID string
var word, wordAdd, meanningAdd string
var counter, counterAddWords int
var pastWords [10]string
var commandsBool = make(map[string]bool)

const (
	botToken   = "botToken"
	baseURL    = "https://heroku_app_name_registered.herokuapp.com/"
	pathToFoto = "static/"
	fileName   = "descriptions2.json"
)

func initTelegram() {
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

func welcomeQuiz(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!"))
}

func newWord(arr []string, fileName string) (string, string) {
	// assign new word, the func doesn't for assign word from rest, rather for smoll words quantity

	w, meanning := "", ""
	for range arr {
		w, meanning = ShowWord(rd.Intn(ShowLength(fileName)), fileName)
		if checkRepetition(arr, w) == false {
			break
		}
	}
	return w, meanning
}

func checkRepetition(arr []string, w string) bool {
	flag := false
	for _, v := range arr {
		if v == w {
			flag = true
			break
		}
	}
	return flag
}

func buttonsOneRightAnswer(shafle bool, str ...string) (int, tgbotapi.InlineKeyboardMarkup) {
	// Creates button and can shafle right answer

	strLenght := len(str)
	// shafle words
	if shafle {
		index := rd.Intn(strLenght)
		str[0], str[index] = str[index], str[0]
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for i := 0; i < strLenght; i += 2 {
		var row []tgbotapi.InlineKeyboardButton
		btn1 := tgbotapi.NewInlineKeyboardButtonData(str[i], str[i]) // word, word
		row = append(row, btn1)
		// if buttons countable and uncountable
		if i <= strLenght {
			btn2 := tgbotapi.NewInlineKeyboardButtonData(str[i+1], str[i+1])
			row = append(row, btn2)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return 0, keyboard
}

func treeRandomWord(fileName, word string, arr []string) (string, string, string) {
	w1, _ := newWord(arr, fileName)
	w2, _ := newWord(arr, fileName)
	w3, _ := newWord(arr, fileName)
	return w1, w2, w3
}

func handleUpdates(update tgbotapi.Update) {
	fmt.Print(bot.Self.UserName)
	// call backes from buttons
	if update.CallbackQuery != nil {
		wordButtonInput := update.CallbackQuery.Data
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Message")
		if counter >= 10 {
			commandsBool["quizBool"] = false
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "10 Quizzz accomplished. Congrates!")
		} else if wordButtonInput == word {
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "That's right «"+word+"»")
		} else {
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Wrong. Right answer is: «"+word+"»")
		}
		bot.Send(msg)
		counter++
		if counter < 10 {
			temp_word, meanning := newWord(pastWords[:], fileName)
			word = temp_word
			pastWords[counter-1] = word
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "\tQuestion №"+strconv.Itoa(counter)+". \nMeanning: "+meanning)
			w1, w2, w3 := treeRandomWord(fileName, word, pastWords[:])
			_, keyboard := buttonsOneRightAnswer(true, word, w1, w2, w3)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Message")
		command := update.Message.Command()
		switch command {

		// typed text
		case "":
			if commandsBool["quizBool"] {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Don't type. Press the button to continue or enter /exit command"))
			} else if commandsBool["addWord"] {
				///////////////////////////////////
				if counterAddWords == 0 {
					wordAdd = update.Message.Text
				} else if counterAddWords == 1 {
					meanningAdd = update.Message.Text
				} else {
					WriteToJsonFile(wordAdd, meanningAdd, "descriptions3.json")
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"For adding word write command /add_word"))
					commandsBool["addWord"] = false
				}
				counterAddWords++
			} else {
				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID,
					pathToFoto+`photo_2021-09-15_23-26-47.jpg`)
				bot.Send(msg)
			}
		case "start":
			//////////////////////////
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "To start 10 word quiz — /start_quizzz\n To exit — /exit")
			bot.Send(msg)
		case "start_quizzz":
			welcomeQuiz(bot, update)
			commandsBool["quizBool"] = true
			counter = 0
			counter++
			temp_word, meanning := newWord(pastWords[:], fileName)
			word = temp_word
			pastWords[counter-1] = word
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\tQuestion №"+strconv.Itoa(counter)+". \nMeanning: "+meanning)
			w1, w2, w3 := treeRandomWord(fileName, word, pastWords[:])
			_, keyboard := buttonsOneRightAnswer(true, word, w1, w2, w3)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		case "add_word":
			counterAddWords = 0
			commandsBool["addWord"] = true
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"Write in 2 messages as such without bracket:\n[word]\n[meanning]"))
		case "exit":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Exit. Done")
			bot.Send(msg)
			for k, _ := range commandsBool {
				commandsBool[k] = false
			}
			// case "add_topic":
		}
	}
}

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	//////// maybe check error
	handleUpdates(update)

	// to monitor changes run: heroku logs --tail
	// log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
}

func main() {
	// Port heroku should give
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// randomazer
	rd.Seed(time.Now().UnixNano())

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	// telegram
	initTelegram()
	commandsBool["quizBool"] = false
	commandsBool["addWord"] = false
	router.POST("/"+bot.Token, webhookHandler)

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
