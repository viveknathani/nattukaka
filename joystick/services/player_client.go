package services

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "joystick/gen/proto"
)

// PlayerClient provides methods to communicate with the player service
type PlayerClient struct {
	client gen.PlayerServiceClient
	conn   *grpc.ClientConn
}

// NewPlayerClient creates a new player client
func NewPlayerClient(playerAddress string) (*PlayerClient, error) {
	conn, err := grpc.NewClient(
		playerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to player service: %w", err)
	}

	client := gen.NewPlayerServiceClient(conn)
	return &PlayerClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (p *PlayerClient) Close() error {
	return p.conn.Close()
}

// BuildImage builds a Docker image on the player node
func (p *PlayerClient) BuildImage(
	ctx context.Context,
	serviceName, repoURL, branch, commitHash string,
) (string, error) {
	req := &gen.BuildImageRequest{
		ServiceName: serviceName,
		RepoUrl:     repoURL,
		Branch:      branch,
		CommitHash:  commitHash,
	}

	resp, err := p.client.BuildImage(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}

	return resp.InternalTag, nil
}

// CreateAndStartContainer creates and starts a container on the player node
func (p *PlayerClient) CreateAndStartContainer(
	ctx context.Context,
	imageTag string,
	envVars map[string]string,
) (string, error) {
	req := &gen.CreateAndStartContainerRequest{
		ImageTag: imageTag,
		Env:      envVars,
	}

	resp, err := p.client.CreateAndStartContainer(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create and start container: %w", err)
	}

	return resp.ContainerId, nil
}

// StopContainer stops a container on the player node
func (p *PlayerClient) StopContainer(ctx context.Context, containerID string) error {
	req := &gen.StopContainerRequest{
		ContainerId: containerID,
	}

	_, err := p.client.StopContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// RemoveContainer removes a container on the player node
func (p *PlayerClient) RemoveContainer(ctx context.Context, containerID string) error {
	req := &gen.RemoveContainerRequest{
		ContainerId: containerID,
	}

	_, err := p.client.RemoveContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}
