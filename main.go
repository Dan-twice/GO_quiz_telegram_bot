package main

import (
	_ "context"
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	database "quiz-bot/db"
)

func main() {
	//init.Heroku()
	//
	//// telegram
	//init.Telegram()
	database.Connection()

	database.Create(database.Dictionary{Word: "testing", Description: "testing word", Context: "context"})

	fmt.Println(database.Select("testing"))
}
