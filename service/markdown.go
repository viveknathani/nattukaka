package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/viveknathani/nattukaka/entity"
)

func (s *Service) GetAllPosts(ctx context.Context, directory string) []entity.Post {

	result := make([]entity.Post, 0)
	path := fmt.Sprintf("static/_md/%s", directory)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		s.Logger.Error(err.Error(), zapReqID(ctx))
	}
	for _, entry := range entries {
		title := entry.Name()
		out, err := exec.Command("git", "log", "-1", "--format=\"%ci\"", "--diff-filter=A", path+title).Output()
		if err != nil {
			s.Logger.Error(err.Error(), zapReqID(ctx))
		}
		output := string(out)
		output = output[1:11]
		result = append(result, entity.Post{
			Title: title,
			Date:  output,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})
	return result
}
