package main

import (
	"context"
	"log"
	"net"
	"player/gen"
	"player/services"

	"google.golang.org/grpc"
)

type playerServer struct {
	gen.UnimplementedPlayerServiceServer
	containerService *services.ContainerService
}

func newPlayerServer() (*playerServer, error) {
	containerService, err := services.NewContainerService()
	if err != nil {
		return nil, err
	}
	return &playerServer{
		containerService: containerService,
	}, nil
}

func (s *playerServer) BuildImage(
	ctx context.Context,
	req *gen.BuildImageRequest,
) (*gen.BuildImageResponse, error) {
	internalTag, err := s.containerService.BuildImage(
		req.ServiceName,
		req.RepoUrl,
		req.Branch,
		req.CommitHash,
	)
	if err != nil {
		return nil, err
	}
	return &gen.BuildImageResponse{
		InternalTag: internalTag,
	}, nil
}

func (s *playerServer) CreateAndStartContainer(
	ctx context.Context,
	req *gen.CreateAndStartContainerRequest,
) (*gen.CreateAndStartContainerResponse, error) {
	containerID, err := s.containerService.CreateAndStartContainer(
		req.ImageTag,
		req.Env,
	)
	if err != nil {
		return nil, err
	}
	return &gen.CreateAndStartContainerResponse{
		ContainerId: containerID,
	}, nil
}

func (s *playerServer) StopContainer(
	ctx context.Context,
	req *gen.StopContainerRequest,
) (*gen.GenericResponse, error) {
	err := s.containerService.StopContainer(
		req.ContainerId,
	)
	if err != nil {
		return nil, err
	}
	return &gen.GenericResponse{
		Message: "Container stopped",
	}, nil
}

func (s *playerServer) RemoveContainer(
	ctx context.Context,
	req *gen.RemoveContainerRequest,
) (*gen.GenericResponse, error) {
	err := s.containerService.RemoveContainer(
		req.ContainerId,
	)
	if err != nil {
		return nil, err
	}
	return &gen.GenericResponse{
		Message: "Container removed",
	}, nil
}

func main() {
	server, err := newPlayerServer()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterPlayerServiceServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	log.Println("gRPC server listening on :50051")
	grpcServer.Serve(listener)
}
