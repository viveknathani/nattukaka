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

var telegramApiKey string = ""

type Message struct {
	ChatID  int64  `json:"chat_id"`
	Caption string `json:"caption"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramApiKey = os.Getenv("TELEGRAM_API_KEY")
}

func getRecords(path string) [][]string {

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
	result = float64(count) / float64(len(health))
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

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 3*vg.Inch, "graph.PNG"); err != nil {
		panic(err)
	}
}

func main() {

	memory := getRecords("/var/memory.txt")
	health := getRecords("/var/health.txt")
	uptime := getUptime(health)
	plotMemoryGraph(memory)

	content := "Hey There!\n"
	content += fmt.Sprintf("uptime: %f", uptime)
	content += "%\n"

	res, err := http.Get("https://hacker-news.firebaseio.com/v0/beststories.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	buffer, bufferError := io.ReadAll(res.Body)

	if bufferError != nil {
		log.Fatal(bufferError)
	}

	var payload interface{}
	json.Unmarshal(buffer, &payload)
	data := payload.([]interface{})

	for i := 0; i < 10; i++ {

		id := data[i].(float64)
		api := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", int64(id))
		res, err := http.Get(api)

		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		buffer, bufferError := io.ReadAll(res.Body)

		if bufferError != nil {
			log.Fatal(bufferError)
		}

		var payload interface{}
		json.Unmarshal(buffer, &payload)
		storyData := payload
		ob := storyData.(map[string]interface{})
		title := ob["title"].(string)
		url := ob["url"].(string)

		content += fmt.Sprintf("%d. %s - %s\n", i+1, title, url)
	}

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

	messageBuf := new(bytes.Buffer)
	messageBufWriter := multipart.NewWriter(messageBuf)
	writerHNWithChat, err := messageBufWriter.CreateFormField("chat_id")
	writerHNWithChat.Write([]byte("5501101308"))
	if err != nil {
		log.Fatal(err)
	}
	writerWithHN, err := messageBufWriter.CreateFormField("text")
	writerWithHN.Write([]byte(content))
	if err != nil {
		log.Fatal(err)
	}
	messageBufWriter.Close()
	_, err = http.Post(telegramHost+"/sendMessage", messageBufWriter.FormDataContentType(), messageBuf)
	if err != nil {
		log.Fatal(err)
	}
}
