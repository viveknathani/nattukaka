package service

import (
	"fmt"
	"log"
	"os"
)

type Service struct{}

func (s *Service) GetAllPosts(directory string) []string {

	result := make([]string, 0)
	path := fmt.Sprintf("static/_md/%s", directory)
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		result = append(result, entry.Name())
	}
	return result
}
