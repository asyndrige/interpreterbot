package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const BASE_URL string = "https://api.telegram.org/bot"

type Conf struct {
	API_KEY string `json:"API_KEY"`
	BOT_URL string `json:"BOT_URL"`
	API_URL string
}

var conf Conf

type TelegramMessage struct {
	UpdateID int `json:"update_id"`
	Message  `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	From      `json:"from"`
	Chat      `json:"chat"`
	Photos    []Photo `json:"photo"`
}

type Photo struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

type From struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func init() {
	file, err := os.Open("conf.json")
	CheckErr(err)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	CheckErr(err)

	conf.API_URL = BASE_URL + conf.API_KEY

	_, err = http.Get(conf.API_URL + "/setWebhook")
	CheckErr(err)
	_, err = http.Get(conf.API_URL + "/setWebhook" + "?url=" + conf.BOT_URL)
	CheckErr(err)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", HttpRestrict(indexHandler, []string{"POST"}))
	http.ListenAndServe(":"+port, nil)
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	var t TelegramMessage
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	CheckErr(err)

	chatId := strconv.Itoa(t.Message.From.ID)
	params := map[string]string{"chat_id": chatId}
	commands := map[string]func([]string) string{
		"/start": start,
	}

	if strings.HasPrefix(t.Message.Text, "/") {

		cmd := strings.Split(t.Message.Text, " ")
		fmt.Println(cmd[0])

		if val, ok := commands[cmd[0]]; ok {
			text := val(cmd[1:])
			params["text"] = text
			go sendMessage(params)
		}

	} else {
		params["text"] = t.Message.From.FirstName + " " + t.Message.From.LastName
		go sendMessage(params)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`{"result": "ok"}`))
}

func sendMessage(getArgs map[string]string) {
	fmt.Println(getArgs)
	apiUrl := conf.API_URL + "/sendMessage?"

	for k, v := range getArgs {
		apiUrl += url.QueryEscape(k)
		apiUrl += "="
		apiUrl += url.QueryEscape(v)
		apiUrl += "&"
	}

	fmt.Println(apiUrl)

	if _, err := http.Get(apiUrl); err != nil {
		log.Println(err)
	}
}

func start(args []string) string {
	return "Starting..."
}

func HttpRestrict(h http.HandlerFunc, verb []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, v := range verb {
			if r.Method != v {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			h(w, r)
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
