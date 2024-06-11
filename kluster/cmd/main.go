package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"ultary.co/kluster/api/grpc/v1"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := grpc.NewServer()
		v1.RegisterSystemServiceServer(server, &System{})

		lis, err := net.Listen("tcp4", "127.0.0.1:50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err = server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}

type System struct {
	v1.SystemServiceServer
}

func (s *System) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}
