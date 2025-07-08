package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

// ContainerService provides methods for container management
type ContainerService struct{}

// NewContainerService creates a new instance of ContainerService
func NewContainerService() *ContainerService {
	return &ContainerService{}
}

// BuildImage builds a new image from the given repository URL and commit hash
func (s *ContainerService) BuildImage(
	serviceName string,
	repoURL string,
	branch string,
	commitHash string,
) (string, error) {
	tag := strings.ToLower(fmt.Sprintf("%s-%s-%d", serviceName, commitHash, time.Now().Unix()))

	baseDir := filepath.Join("./temp/repo", tag)
	absDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(absDir, 0755); err != nil {
		return "", err
	}

	repo, err := git.PlainClone(absDir, &git.CloneOptions{
		URL:           repoURL,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})

	if err != nil {
		fmt.Println("error cloning repository: " + err.Error())
		return "", err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("error checking out commit: " + err.Error())
		return "", err
	}
	worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})

	command := exec.Command("docker", "build", "-t", tag, ".")
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Dir = absDir

	if err := command.Run(); err != nil {
		fmt.Println("error building image: " + err.Error())
		return "", err
	}

	return tag, nil
}
