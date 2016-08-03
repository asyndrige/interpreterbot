package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Main Bot functionality.

type Commander interface {
	start([]string) string
	help([]string) string
}

type Bot struct{}

type command func([]string) string

var commands = make(map[string]command)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	var t TelegramMessage
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	checkErr(err)

	result := make(chan bool)
	chatID := strconv.Itoa(t.Message.From.ID)
	params := map[string]string{"chat_id": chatID}

	if strings.HasPrefix(t.Message.Text, "/") {
		go execCommand(params, t.Message.Text, result)
	} else {
		params["text"] = fmt.Sprintf("You are %s %s.\n Type /help to see list of available commands.", t.Message.From.FirstName, t.Message.From.LastName)
		go sendMessage(params, result)
	}

	res.Header().Set("Content-Type", "application/json")
	if sendRes := <-result; sendRes {
		res.Write([]byte(`{"result": "ok"}`))
	} else {
		res.Write([]byte(`{"result": 0}`))
	}
}

func sendMessage(getArgs map[string]string, res chan bool) {
	fmt.Println(getArgs)
	apiURL := conf.API_URL + "/sendMessage?"

	for k, v := range getArgs {
		apiURL += url.QueryEscape(k)
		apiURL += "="
		apiURL += url.QueryEscape(v)
		apiURL += "&"
	}

	fmt.Println(apiURL)

	if _, err := http.Get(apiURL); err != nil {
		log.Println(err)
		res <- false
	} else {
		res <- true
	}
}

func execCommand(getArgs map[string]string, text string, result chan bool) {
	commands["/start"] = start
	commands["/help"] = help

	cmd := strings.Split(text, " ")
	fmt.Println(cmd[0])

	if val, ok := commands[cmd[0]]; ok {
		text := val(cmd[1:])
		getArgs["text"] = text
		go sendMessage(getArgs, result)
	} else {
		text := "Command is not supported.\nType /help to see list of commands."
		getArgs["text"] = text
		go sendMessage(getArgs, result)
	}
}

func start(args []string) string {
	return "Starting..."
}

func help(args []string) string {
	res := "Available commands:\n"
	for k := range commands {
		res += k
		res += "\n"
	}
	return res
}

