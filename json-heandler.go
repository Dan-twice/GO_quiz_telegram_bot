package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Word struct {
	Word     string `json:"word"`
	Meanning string `json:"description"`
}

// we initialize
var words []Word

func WriteToJsonFile(word, meanning, fileName string) {
	data := Word{
		Word:     word,
		Meanning: meanning,
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(fileName, file, 0644)
}

func ShowWord(index int, fileName string) (string, string) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	// func Unmarshal(data []byte, v interface{}) error
	json.Unmarshal(byteValue, &words)

	return words[index].Word, words[index].Meanning
	// return words.Wds[index].Word, words.Wds[index].Meanning
}

func ShowLength(fileName string) int {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &words)
	return len(words)
}
