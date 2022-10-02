package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const MAX_HN_STORIES = 10

var telegramApiKey string = ""

type HackerNewsAPIBody struct {
	ChatId         int64  `json:"chat_id"`
	Text           string `json:"text"`
	DisablePreview bool   `json:"disable_web_page_preview"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramApiKey = os.Getenv("TELEGRAM_API_KEY")
}

func getRecordsFromCSVFile(path string) [][]string {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func getUptime(health [][]string) float64 {

	result := 0.0
	count := 0
	for _, dataPoint := range health {
		if dataPoint[1] == "200" {
			count++
		}
	}
	result = float64(count) * 100.0 / float64(len(health))
	return result
}

func plotMemoryGraph() {

	cmd := exec.Command("python3", "scripts/graph.py")
	err := cmd.Run()
	//cmd.Stderr = os.Stderr
	if err != nil {
		log.Fatal(err)
	}
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

func sendToTelegram(content string, photoPath string) {

	telegramHost := fmt.Sprintf("https://api.telegram.org/bot%s", telegramApiKey)
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	file, err := os.Open("graph.PNG")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writerWithPhoto, err := writer.CreateFormFile("photo", "graph")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(writerWithPhoto, file)
	writerWithChat, err := writer.CreateFormField("chat_id")
	writerWithChat.Write([]byte("5501101308"))
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()
	_, err = http.Post(telegramHost+"/sendPhoto", writer.FormDataContentType(), buf)
	if err != nil {
		log.Fatal(err)
	}
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

	health := getRecordsFromCSVFile("/var/health.txt")
	uptime := getUptime(health)
	plotMemoryGraph()
	data := getTopStoriesFromHN()

	content := "Hey There!\n"
	content += fmt.Sprintf("uptime: %f", uptime)
	content += "%\n"
	for i := 0; i < MAX_HN_STORIES; i++ {
		id := int64(data[i].(float64))
		data := getHNStory(id)
		ob := data.(map[string]interface{})
		title := ob["title"].(string)
		url := fmt.Sprintf("news.ycombinator.com/item?id=%d", id)
		content += fmt.Sprintf("%d. %s - %s\n", i+1, title, url)
	}

	sendToTelegram(content, "graph.PNG")
}
