package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	//Framework
	models "github.com/mccraymt/ms-black-history/models"
	"github.com/gorilla/mux"
)

// HandlePostalCodeLookup returns rating data for a postal code
func HandlePostalCodeLookup(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	postalCode := args["code"]
	//fmt.Fprintf(w, "Code: %v\n", postalCode)

	pcRecord, _ := models.PostalCodeDict[postalCode]
	if pcRecord == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		pcRecord.ID = postalCode
		w.WriteHeader(http.StatusOK)
		foo, _ := json.MarshalIndent(pcRecord, "", "  ")
		bar := string(foo)

		fmt.Fprintf(w, bar)
	}
}
