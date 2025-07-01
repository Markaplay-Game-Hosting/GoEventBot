package main

import (
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"net/http"
)

func (app *application) referenceHandler(w http.ResponseWriter, r *http.Request) {
	var scheme string
	if r.TLS != nil {
		scheme = "https:"
	} else {
		scheme = "http:"
	}
	specUrl := fmt.Sprintf("%s//%s/docs/swagger.json", scheme, r.Host)
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
		SpecURL:  specUrl,
		DarkMode: true,
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Go Event Bot API",
		},
	})

	if err != nil {
		fmt.Printf("%v", err)
	}

	_, err = fmt.Fprintln(w, htmlContent)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
