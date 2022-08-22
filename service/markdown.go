package service

import (
	"context"
	"fmt"
	"os"
)

func (s *Service) GetAllPosts(ctx context.Context, directory string) []string {

	result := make([]string, 0)
	path := fmt.Sprintf("static/_md/%s", directory)
	entries, err := os.ReadDir(path)
	if err != nil {
		s.Logger.Error(err.Error(), zapReqID(ctx))
	}
	for _, entry := range entries {
		result = append(result, entry.Name())
	}
	return result
}
