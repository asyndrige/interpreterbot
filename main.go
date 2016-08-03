package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type Conf struct {
	API_KEY string `json:"API_KEY"`
	BOT_URL string `json:"BOT_URL"`
	API_URL string
}

var conf Conf

const BASE_URL string = "https://api.telegram.org/bot"

func init() {

	file, err := os.Open("conf.json")
	checkErr(err)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	checkErr(err)

	conf.API_URL = BASE_URL + conf.API_KEY

	_, err = http.Get(conf.API_URL + "/setWebhook")
	checkErr(err)
	_, err = http.Get(conf.API_URL + "/setWebhook" + "?url=" + conf.BOT_URL)
	checkErr(err)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", httpRestrict(indexHandler, []string{"POST"}))
	http.ListenAndServe(":"+port, nil)
}
