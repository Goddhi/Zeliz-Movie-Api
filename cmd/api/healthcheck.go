package main

import (

	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": 	"available",
		"system_info": map[string]string{
			"enviroment": app.config.env,
			"version": 	version, // remember version is a global vaiable
		},

	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}


}
