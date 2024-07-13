package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ultary/monokube/kluster/api/grpc/v1"
)

func (s *Server) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}

func (s *Server) EnableOpenTelemetryCollector(context.Context, *v1.EnableOpenTelemetryCollectorRequest) (*v1.EnableOpenTelemetryCollectorResponse, error) {
	s.cluster.System().EnableOpenTelemetryCollector()
	return &v1.EnableOpenTelemetryCollectorResponse{}, nil
}

func (s *Server) DisableOpenTelemetryCollector(context.Context, *v1.DisableOpenTelemetryCollectorRequest) (*v1.DisableOpenTelemetryCollectorResponse, error) {
	s.cluster.System().DisableOpenTelemetryCollector()
	return &v1.DisableOpenTelemetryCollectorResponse{}, nil
}

func (s *Server) UpdateOpenTelemetryCollector(context.Context, *v1.UpdateOpenTelemetryCollectorRequest) (*v1.UpdateOpenTelemetryCollectorResponse, error) {
	s.cluster.System().UpdateOpenTelemetryCollector()
	return &v1.UpdateOpenTelemetryCollectorResponse{}, nil
}
