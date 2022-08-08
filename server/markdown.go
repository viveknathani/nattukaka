package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

type blogIndexPageVariables struct {
	PostList []string
}

func (s *Server) serveMarkdownIndex(w http.ResponseWriter, r *http.Request) {

	directory := strings.TrimPrefix(r.URL.Path, "/")
	indexFilePath := "static/pages/posts.html"
	t, err := template.ParseFiles(indexFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	err = t.Execute(w, blogIndexPageVariables{
		PostList: s.Service.GetAllPosts(directory),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

func addViewPort(input []byte) []byte {

	const keyword = "</title>"
	stream := string(input)
	idx := strings.Index(stream, keyword)
	if idx == -1 {
		log.Fatal("Failed to add meta viewport")
	}
	end := idx + len(keyword) - 1
	arr := make([]byte, 0)
	for i := 0; i <= end; i++ {
		arr = append(arr, byte(stream[i]))
	}

	content := `<meta name="viewport" content="width=device-width, initial-scale=1">`

	for i := 0; i < len(content); i++ {
		arr = append(arr, byte(content[i]))
	}

	for i := end + 1; i < len(stream); i++ {
		arr = append(arr, byte(stream[i]))
	}

	return arr
}

func markdowntoHTML(source string, title string, cssFile string) []byte {

	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stream, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	params := blackfriday.HTMLRendererParameters{
		CSS:   cssFile,
		Title: title,
		Flags: blackfriday.CompletePage,
	}
	renderer := blackfriday.NewHTMLRenderer(params)
	return addViewPort(blackfriday.Run(stream, blackfriday.WithRenderer(renderer)))
}

func (s *Server) serveMarkdownPost(w http.ResponseWriter, r *http.Request) {

	postPath := fmt.Sprintf("static/_md%s", r.URL.Path)
	title := mux.Vars(r)["title"]
	html := markdowntoHTML(postPath, title, "/static/styles/theme.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}
