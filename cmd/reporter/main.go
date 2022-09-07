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
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const MAX_HN_STORIES = 10

var telegramApiKey string = ""

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

func plotMemoryGraph(memory [][]string) {

	timeStamps := make([]float64, 0)
	for _, dataPoint := range memory {
		seconds, err := strconv.ParseInt(dataPoint[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		timeStamps = append(timeStamps, float64(seconds))
	}

	megabytes := make([]float64, 0)
	for _, dataPoint := range memory {
		megabyte, err := strconv.ParseFloat(strings.TrimSpace(dataPoint[1]), 64)
		if err != nil {
			log.Fatal(err)
		}
		megabytes = append(megabytes, megabyte/1000)
	}

	pts := make(plotter.XYs, len(megabytes))
	for i := range pts {
		pts[i].X = timeStamps[i]
		pts[i].Y = megabytes[i]
	}

	p := plot.New()

	p.Title.Text = "nattukaka - memory usage"
	p.X.Label.Text = "timestamp (unix epoch)"
	p.Y.Label.Text = "memory (MB)"

	err := plotutil.AddLinePoints(p, pts)
	if err != nil {
		panic(err)
	}

	if err := p.Save(4*vg.Inch, 3*vg.Inch, "graph.PNG"); err != nil {
		panic(err)
	}
}

func getTopStoriesFromHN() []interface{} {

	response, err := http.Get("https://hacker-news.firebaseio.com/v0/beststories.json")
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
	buf = new(bytes.Buffer)
	writer = multipart.NewWriter(buf)
	writerWithChat, err = writer.CreateFormField("chat_id")
	writerWithChat.Write([]byte("5501101308"))
	if err != nil {
		log.Fatal(err)
	}
	writerWithContent, err := writer.CreateFormField("text")
	writerWithContent.Write([]byte(content))
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()
	_, err = http.Post(telegramHost+"/sendMessage", writer.FormDataContentType(), buf)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	memory := getRecordsFromCSVFile("/var/memory.txt")
	health := getRecordsFromCSVFile("/var/health.txt")
	uptime := getUptime(health)
	plotMemoryGraph(memory)
	data := getTopStoriesFromHN()

	content := "Hey There!\n"
	content += fmt.Sprintf("uptime: %f", uptime)
	content += "%\n"
	for i := 0; i < MAX_HN_STORIES; i++ {
		id := data[i].(float64)
		data := getHNStory(int64(id))
		ob := data.(map[string]interface{})
		title := ob["title"].(string)
		url := ob["url"].(string)
		content += fmt.Sprintf("%d. %s - %s\n", i+1, title, url)
	}

	sendToTelegram(content, "graph.PNG")
}
