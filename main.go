package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Conf struct {
	API_KEY        string `json:"API_KEY"`
	BOT_URL        string `json:"BOT_URL"`
	API_URL        string
	DETECT_API_KEY string `json: DETECT_API_KEY`
	DETECT_API_URL string `json: DETECT_API_URL`
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

	urlQuery := url.Values{}
	urlQuery.Add("access_key", conf.DETECT_API_KEY)
	_, err = http.Get(conf.DETECT_API_URL + "detect?" + urlQuery.Encode())
	fmt.Println(conf.DETECT_API_URL + "detect?" + urlQuery.Encode())
	checkErr(err)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", httpRestrict(indexHandler, []string{"POST"}))
	http.ListenAndServe(":"+port, nil)
}
