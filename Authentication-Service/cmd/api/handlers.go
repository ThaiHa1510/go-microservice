package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for authentication")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println("err read data", err)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		app.LogItem(requestPayload.Email, requestPayload.Password, false)
		return
	}

	valid := user.MachePassword(requestPayload.Password)
	if !valid {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		app.LogItem(requestPayload.Email, requestPayload.Password, false)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    user,
	}
	go app.LogItem(requestPayload.Email, requestPayload.Password, true)
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogItem(email string, password string, isError bool) bool {
	type LogData struct {
		Email    string
		Password string
		Error    bool
	}
	var requestPayload struct {
		Name string `json:"name"`
		Data LogData
	}
	requestPayload.Name = "authentication"
	requestPayload.Data.Email = email
	requestPayload.Data.Password = password
	requestPayload.Data.Error = isError
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	logServiceURL := "http://logger-service/log"
	response, err := http.Post(logServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("err save log", err)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Println("err save log", err)
		return false
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		log.Println("err save log", err)
		return false
	}
	log.Println("save log success")
	return true

}
