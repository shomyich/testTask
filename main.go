package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var baseURL = getEnv("BASE_URL", "https://api.green-api.com/")

type Message struct {
	ChatId  string `json:"chatId"`
	Message string `json:"message"`
}

type File struct {
	ChatId   string `json:"chatId"`
	UrlFile  string `json:"urlFile"`
	FileName string `json:"fileName"`
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/getSettings", getSettings)
	http.HandleFunc("/getStateInstance", getStateInstance)
	http.HandleFunc("/sendMessage", sendMessage)
	http.HandleFunc("/sendFileByUrl", sendFileByUrl)

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")

	url := baseURL + "waInstance" + idInstance + "/getSettings/" + apiTokenInstance
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to get settings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get settings: "+response.Status, response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	w.Write(body)
}

func getStateInstance(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")

	url := baseURL + "waInstance" + idInstance + "/getStateInstance/" + apiTokenInstance
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to get state instance: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get state instance: "+response.Status, response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	w.Write(body)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	chatId := r.URL.Query().Get("chatId")
	message := r.URL.Query().Get("message")

	msg := Message{
		ChatId:  chatId,
		Message: message,
	}

	url := baseURL + "waInstance" + idInstance + "/sendMessage/" + apiTokenInstance
	jsonValue, _ := json.Marshal(msg)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to send message: "+response.Status, response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	w.Write(body)
}

func sendFileByUrl(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	chatId := r.URL.Query().Get("chatId")
	urlFile := r.URL.Query().Get("fileUrl")
	fileName := r.URL.Query().Get("fileName")

	file := File{
		ChatId:   chatId,
		UrlFile:  urlFile,
		FileName: fileName,
	}

	url := baseURL + "waInstance" + idInstance + "/sendFileByUrl/" + apiTokenInstance
	jsonValue, _ := json.Marshal(file)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "Failed to send file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to send file: "+response.Status, response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	w.Write(body)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
