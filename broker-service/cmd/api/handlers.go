package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload)
	case "log":
		//app.logItem(w, requestPayload.Log)
		app.LogItemViaRpc(w, requestPayload.Log)
	case "mail":
		app.SendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) logEventViaRabit(w http.ResponseWriter, entry LogPayload) {
	fmt.Println("logEvent received")
	err := app.pushToQueue(entry.Name, entry.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Log via RabbitMQ"),
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emmitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	j, _ := json.MarshalIndent(payload, "", "\t")
	err = emmitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	fmt.Println("logItem received")
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		fmt.Println("Error create request")
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		fmt.Println("Error sending request", err.Error())
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling logger service"))
		fmt.Println("error calling logger service")
		return
	}
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		fmt.Println("Error decoding response from logger service")
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err)
		fmt.Println("Error jsonFromService")
		return
	}
	var result jsonResponse
	result.Message = "logged"
	result.Data = jsonFromService.Data
	result.Error = false
	app.writeJSON(w, http.StatusAccepted, result)
}
func (app *Config) authenticate(w http.ResponseWriter, payload RequestPayload) {
	fmt.Println("authenticating email", payload.Auth.Email)
	jsonData, _ := json.MarshalIndent(payload.Auth, "", "\t")
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer([]byte(jsonData)))

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return

	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var result jsonResponse
	result.Message = "Authenticated"
	result.Data = jsonFromService.Data
	result.Error = false
	app.writeJSON(w, http.StatusAccepted, result)
	return

}

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")
	// call mail service
	fmt.Println("sending mail to", msg.To)
	mailServierURL := "http://mail-service/send"
	request, err := http.NewRequest("POST", mailServierURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To
	app.writeJSON(w, http.StatusAccepted, payload)

}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) LogItemViaRpc(w http.ResponseWriter, payload LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayLoad := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}

	defer client.Close()

	var result string
	err = client.Call("RPCServer.LogItem", rpcPayLoad, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.writeJSON(w, http.StatusAccepted, resp)

}
