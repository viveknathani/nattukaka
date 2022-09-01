package entity

type Note struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
