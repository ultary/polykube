package v1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ultary/monokube/core/api/grpc/v1"
)

func RegisterSystemServer(s grpc.ServiceRegistrar) {
	v1.RegisterSystemServer(s, &System{})
}

type System struct {
	v1.SystemServer
}

func (s *System) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}
