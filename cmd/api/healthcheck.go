package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envolope{
		"status": "available",
		"system_info": map[string]string{
			"envirnoment": app.config.env,
			"version":     version,
		},
	}

	// Marshal function takes an input of interface{}
	// and send out a bytestream in json format
	err := app.writeJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
