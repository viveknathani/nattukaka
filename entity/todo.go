package entity

type Todo struct {
	Id          string  `json:"id"`
	UserId      string  `json:"userId"`
	Task        string  `json:"task"`
	Status      string  `json:"status"`
	Deadline    *string `json:"deadline"`
	CompletedAt *string `json:"completedAt"`
}
