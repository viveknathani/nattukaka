package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

// ContainerService provides methods for container management
type ContainerService struct {
	dockerClient *client.Client
}

// NewContainerService creates a new instance of ContainerService
func NewContainerService() (*ContainerService, error) {
	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		fmt.Println("error creating docker client: " + err.Error())
		return nil, err
	}
	return &ContainerService{
		dockerClient: dockerClient,
	}, nil
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

func envMapToSlice(m map[string]string) []string {
	env := make([]string, 0, len(m))
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

// CreateAndStartContainer creates and starts a container from the given image
func (s *ContainerService) CreateAndStartContainer(
	tag string,
	environmentVariables map[string]string,
) error {
	fmt.Println("Starting container with tag: " + tag)

	fmt.Println("Container create started")
	containerCreateResponse, err := s.dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: tag,
			Env:   envMapToSlice(environmentVariables),
		},
		&container.HostConfig{
			NetworkMode: "host",
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		fmt.Println("error creating container: " + err.Error())
		return err
	}
	fmt.Println("Container created with ID: " + containerCreateResponse.ID)

	err = s.dockerClient.ContainerStart(
		context.Background(),
		containerCreateResponse.ID,
		container.StartOptions{},
	)
	if err != nil {
		fmt.Println("error starting container: " + err.Error())
		return err
	}
	fmt.Println("Container started")
	return nil
}

// StopContainer stops a container by ID
func (s *ContainerService) StopContainer(
	containerID string,
) error {
	s.dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{})
	return nil
}

// RemoveContainer removes a container by ID
func (s *ContainerService) RemoveContainer(
	containerID string,
) error {
	s.dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})
	return nil
}
