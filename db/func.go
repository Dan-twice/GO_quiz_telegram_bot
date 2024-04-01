package database

import (
	"fmt"
	"log"
)

func Create(dict Dictionary) {
	tx := db.MustBegin()
	fmt.Println("here 0")

	tx.MustExec("INSERT INTO dictionary (word, description, context) VALUES (:word, :description, :context)",
		dict)

	fmt.Println("here 1")

	err := tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func Update() {

}

func Delete() {

}

func Select(w string) *Dictionary {
	dict := []Dictionary{}

	fmt.Println("here 2")

	err := db.Select(&dict, "SELECT * FROM dictionary WHERE word=:w", map[string]interface{}{"w": w})
	if err != nil {
		log.Println(err)
	}

	return &dict[0]
}
