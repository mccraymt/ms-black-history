package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// FlashCard ...
type FlashCard struct {
	ID         string   `json:"id"`
	Volume     string   `json:"value"`
	Title      string   `json:"title"`
	Subtitle   string   `json:"subtitle"`
	Facts      []string `json:"facts"`
	ImageSrc   string   `json:"imageSrc"`
	References []string `json:"references"`
}

// FlashCardDict ...
var FlashCardDict map[string]*FlashCard

// FlashCards ...
var FlashCards []*FlashCard

func init() {
	FlashCardDict = make(map[string]*FlashCard)
	flashCards, err := getFlashCards()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, r := range flashCards {
		if r.ID != "" {
			FlashCardDict[r.ID] = r
		}
	}

	FlashCards = flashCards

}

func getFlashCards() ([]*FlashCard, error) {
	raw, err := ioutil.ReadFile("./models/cards.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []*FlashCard
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return c, nil
}
