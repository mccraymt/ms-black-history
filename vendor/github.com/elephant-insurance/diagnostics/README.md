# Elephant diagnostics library
A library for retrieving and serving diagnostic data from web microservices.

## How do I hook this up?
Include the library:
```include "github.com/elephant-insurance/diagnostics"```

After that, it only takes two lines. First initialize the package and get back an `http.Handler`:
```
diagnosticsHandler, err := diagnostics.Initialize(config.Config.Version, logrus.New())
```
The logger argument can be anything that implements `Printf`. Set it to nil to suppress logging of results.
Now, if the service is using Negroni and gorilla/mux, simply add the handler to a route (by default we use GET "/"):
```
mainRouter.Methods("GET").Path("/").Handler(diagnosticsHandler)
```
If the service is still running Martini, you'll need to be a bit more explicit:
```
m.Get("/", func(w http.ResponseWriter, r *http.Request) { diagnosticsHandler.ServeHTTP(w, r) })
```
Now when this route is called, you will see an extensive set of diagnostic information about the running process. 

## How do I test connectivity to other services that I need?
If you want to check for connectivity to other services when diagnostics are run, you can add a test with a single line:
```
diagnostics.AddConnectionTest("http://my-other-service.elephant.com")
```

Now, every time the diagnostics page is loaded with a query string parameter `runTests=true`, your connection test will run and display its results.

## How do I run other kinds of tests?
If you want to run more complex tests when the diagnostic page is loaded, create one or more functions with this signature:
```
func myDiagnostic() (diagnostics.DiagnosticResult, error) {
	rp := diagnostics.NewResult().Succeed().SetDescriptionf("Data dictionary has %v entries", len(models.PostalCodeDict))

	return *rp, nil
}
```

Note that you can use this "fluent" syntax to create your return value, or you can just create a `DiagnosticResult` struct and set its properties yourself.

Now just add it to the list of tests to run:
`diagnostics.AddDiagnosticTest("Diagnostic test example", myDiagnostic)`

## What does the diagnostic page look like?
The diagnostic page returns a result of type `SystemInfo` marshaled into JSON. If you want to read this output in another process, you can import this type into other projects and deserialize the diagnostic body:
```
	si := SystemInfo{}
	err = json.Unmarshal(response.Body.Bytes(), &si)
	if err != nil {
		fmt.Printf("Failed to unmarshal response: %v \n", err.Error())
	}
```
Now you have an in-memory `SystemInfo` whose properties you can use.

If any of the tests you add fail (e.g., `return diagnostics.NewResult().Fail().SetDescription("A bad thing happened.")`), the diagnostics page will return with a 500 code. If all tests succeed or leave the Success field empty, the page will return with 200 OK.
