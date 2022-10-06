package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const MAX_HN_STORIES = 10

var telegramApiKey string = ""

type HackerNewsAPIBody struct {
	ChatId         int64  `json:"chat_id"`
	Text           string `json:"text"`
	DisablePreview bool   `json:"disable_web_page_preview"`
}

func init() {
	telegramApiKey = os.Getenv("TELEGRAM_API_KEY")
}

func getTopStoriesFromHN() []interface{} {

	response, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	buffer, bufferError := io.ReadAll(response.Body)
	if bufferError != nil {
		log.Fatal(bufferError)
	}
	var payload interface{}
	json.Unmarshal(buffer, &payload)
	data := payload.([]interface{})
	return data
}

func getHNStory(id int64) interface{} {

	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	buffer, bufferError := io.ReadAll(response.Body)
	if bufferError != nil {
		log.Fatal(bufferError)
	}
	var payload interface{}
	json.Unmarshal(buffer, &payload)
	return payload
}

func sendToTelegram(content string) {

	telegramHost := fmt.Sprintf("https://api.telegram.org/bot%s", telegramApiKey)
	body := &HackerNewsAPIBody{
		ChatId:         5501101308,
		Text:           content,
		DisablePreview: true,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.Post(telegramHost+"/sendMessage", "application/json", bytes.NewReader(bodyJSON))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	data := getTopStoriesFromHN()
	content := ""
	for i := 0; i < MAX_HN_STORIES; i++ {
		id := int64(data[i].(float64))
		data := getHNStory(id)
		ob := data.(map[string]interface{})
		title := ob["title"].(string)
		url := fmt.Sprintf("news.ycombinator.com/item?id=%d", id)
		content += fmt.Sprintf("%d. %s - %s\n", i+1, title, url)
	}

	sendToTelegram(content)
}
