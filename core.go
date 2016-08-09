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

type LanguageLayerMessage struct {
	Success bool `json:"success"`
	Results []struct {
		LanguageCode   string  `json:"language_code"`
		LanguageName   string  `json:"language_name"`
		Probability    float64 `json:"probability"`
		Percentage     float64 `json:"percentage"`
		ReliableResult bool    `json:"reliable_result"`
	} `json:"results"`
}

type Command struct {
	executor     func(url.Values, []string) string
	argsRequired bool
}

type Phrase struct {
	Content           string
	DetectedLang      string
	TranslatedContent string
}

func (p *Phrase) translate() string {
	p.TranslatedContent = "Translating"
	return p.TranslatedContent
}

var (
	phrase       Phrase
	commands     = make(map[string]Command)
	requiredArgs = make(map[string]bool)
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	var t TelegramMessage
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	checkErr(err)

	result := make(chan bool)
	chatID := strconv.Itoa(t.Message.From.ID)
	params := url.Values{}
	params.Add("chat_id", chatID)

	if strings.HasPrefix(t.Message.Text, "/") {
		go execCommand(params, t.Message.Text, result)
	} else {
		params.Add("text", fmt.Sprintf("You are %s %s.\n Type /help to see list of available commands.", t.Message.From.FirstName, t.Message.From.LastName))
		go sendMessage(params, result)
	}

	res.Header().Set("Content-Type", "application/json")
	if sendRes := <-result; sendRes {
		res.Write([]byte(`{"result": "ok"}`))
	} else {
		res.Write([]byte(`{"result": 0}`))
	}
}

func sendMessage(getArgs url.Values, res chan<- bool) {
	fmt.Println(getArgs)
	apiURL := conf.API_URL + "/sendMessage?"
	apiURL += getArgs.Encode()

	if _, err := http.Get(apiURL); err != nil {
		log.Println(err)
		res <- false
	} else {
		res <- true
	}
}

func execCommand(getArgs url.Values, text string, result chan<- bool) {
	commands["/start"] = Command{start, false}
	commands["/help"] = Command{help, false}
	commands["/detectlang"] = Command{detectLang, true}

	cmd := strings.Split(text, " ")
	fmt.Println(cmd[0])

	if val, ok := commands[cmd[0]]; ok {
		var text string

		if val.argsRequired {
			if len(cmd[1:]) != 0 {
				text = val.executor(getArgs, cmd[1:])
			} else {
				text = "This command requires at least 1 argument."
			}
		} else {
			text = val.executor(getArgs, cmd[1:])
		}

		getArgs.Add("text", text)
		go sendMessage(getArgs, result)
	} else {
		text := "Command is not supported.\nType /help to see list of commands."
		getArgs.Add("text", text)
		go sendMessage(getArgs, result)
	}
}

func start(_ url.Values, args []string) string {
	return "Starting..."
}

func help(_ url.Values, args []string) string {
	res := "Available commands:\n"
	for k := range commands {
		res += k
		res += "\n"
	}
	return res
}

func detectLang(_ url.Values, args []string) string {
	phrase.Content = strings.Join(args, " ")
	detectAPIURL := conf.DETECT_API_URL + "detect?"
	urlQuery := url.Values{}
	urlQuery.Add("access_key", conf.DETECT_API_KEY)
	urlQuery.Add("query", phrase.Content)

	res, err := http.Get(detectAPIURL + urlQuery.Encode())
	fmt.Println(detectAPIURL + urlQuery.Encode())
	checkErr(err)

	var l LanguageLayerMessage
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&l)
	checkErr(err)

	phrase.DetectedLang = l.Results[0].LanguageName

	resText := []string{
		"Detected language is ",
		l.Results[0].LanguageName,
		" with the probability of ",
		strconv.Itoa(int(l.Results[0].Percentage)),
		"%.\n"}

	if len(l.Results) > 1 {
		resText = append(resText, "Others are:\n")
		for _, v := range l.Results[1:] {
			resText = append(resText, v.LanguageName+": "+strconv.Itoa(int(v.Percentage))+"%.\n")
		}
	}

	return strings.Join(resText, "")
}
