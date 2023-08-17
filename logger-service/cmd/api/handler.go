package main

import (
	"fmt"
	"log"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for logging")
	//read json into var
	var requesPayload JSONPayload
	_ = app.readJSON(w, r, &requesPayload)

	// insert data
	event := data.LogEntry{
		Name: requesPayload.Name,
		Data: requesPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	app.writeJSON(w, http.StatusAccepted, resp)

}
