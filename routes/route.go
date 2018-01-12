package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	//Framework
	//c "../../config"

	//Controllers
	c "github.com/mccraymt/ms-black-history/controllers"
)

// New initializes routes for the app
func New(diagnosticsHandler http.Handler) *mux.Router {

	//Create main router
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.KeepContext = true

	//App Routes
	mainRouter.Methods("GET").Path("/").Handler(diagnosticsHandler)
	mainRouter.Methods("GET").Path("/error").HandlerFunc(c.SayError)
	mainRouter.Methods("GET").Path("/postal-codes/{code}").HandlerFunc(c.HandlePostalCodeLookup)

	mainRouter.Methods("GET").Path("/flash-cards").HandlerFunc(c.HandleFlashCardListAll)
	mainRouter.Methods("GET").Path("/flash-cards/{index}").HandlerFunc(c.HandleFlashCardLookup)

	return mainRouter

}
