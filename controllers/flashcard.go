package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	//Framework
	"github.com/gorilla/mux"
	models "github.com/mccraymt/ms-black-history/models"
)

// HandleFlashCardListAll returns rating data for a postal code
func HandleFlashCardListAll(w http.ResponseWriter, r *http.Request) {
	record := models.FlashCards
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	foo, _ := json.MarshalIndent(record, "", "  ")
	bar := string(foo)

	fmt.Fprintf(w, bar)
}

// HandleFlashCardLookup returns rating data for a postal code
func HandleFlashCardLookup(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	index := args["index"]
	//fmt.Fprintf(w, "Code: %v\n", postalCode)

	pcRecord, _ := models.FlashCardDict[index]
	// pcRecord := models.FlashCardDict[index]
	if pcRecord == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		pcRecord.ID = index
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		foo, _ := json.MarshalIndent(pcRecord, "", "  ")
		bar := string(foo)

		fmt.Fprintf(w, bar)
	}
}
