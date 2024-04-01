package handler

import (
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"strconv"
)

var userID string
var word, wordAdd, meanningAdd string
var counter, counterAddWords int
var pastWords [10]string
var commandsBool = make(map[string]bool)

//// randomazer
//rd.Seed(time.Now().UnixNano())

const (
	pathToFoto = "static/"
)

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

func welcomeQuiz(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!"))
}
