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

// GeAllPosts will read through a directory of posts
// and list out all the entries, sorted by their
// first commit date in descending order.
func (service *Service) GetAllPosts(ctx context.Context, directory string) ([]entity.Post, error) {

	result := make([]entity.Post, 0)

	path := fmt.Sprintf("static/_md/%s", directory)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoPostList
	}

	postType := strings.TrimSuffix(strings.TrimPrefix(path, "static/_md/"), "/")
	for _, entry := range entries {

		title := entry.Name()
		out, err := exec.Command("git", "log", "-1", "--format=\"%ci\"", "--diff-filter=A", path+title).Output()
		if err != nil {
			service.Logger.Error(err.Error(), zapReqID(ctx))
			return nil, ErrNoPostList
		}
		output := string(out)
		output = output[1:11] // output is `"yyyy-mm-dd..."`, we want just `yyyy-mm-dd`
		result = append(result, entity.Post{
			Type:  postType,
			Title: title,
			Date:  output,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})
	return result, nil
}
