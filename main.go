package main

import (
	"fmt"
	"net/http"

	logrus "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/elephant-insurance/diagnostics"

	"github.com/mccraymt/ms-black-history/config"
	"github.com/mccraymt/ms-black-history/models"
	routes "github.com/mccraymt/ms-black-history/routes"
	"github.com/rs/cors"
)

func main() {
	diagnosticsHandler, err := diagnostics.Initialize(config.Config.Version, logrus.New())
	if err != nil {
		logrus.Panic("Could not initialize diagnostics:\n" + err.Error())
	}

	diagnostics.AddDiagnosticTest("Trivial diagnostic result example", trivialDiagnostic)
	// This is an example of how to add a dependency check with one line:
	// diagnostics.AddConnectionTest("http://google.com")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(c)
	http.Handle("/", app)
	app.UseHandler(routes.New(diagnosticsHandler))
	port := fmt.Sprintf(":%v", config.Config.Port)
	app.Run(port)
}

func trivialDiagnostic() (diagnostics.DiagnosticResult, error) {
	rp := diagnostics.NewResult().Succeed().SetDescriptionf("Data dictionary has %v entries", len(models.PostalCodeDict))

	return *rp, nil
}
