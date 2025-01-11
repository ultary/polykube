package system

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ultary/polykube/kluster/api/grpc/v1"
	"github.com/ultary/polykube/kluster/pkg/k8s"
	"github.com/ultary/polykube/kluster/pkg/kube/system/apps/otlp"
)

//var _ v1.SystemServer = Server{}

type Server struct {
	v1.UnimplementedSystemServer
	cluster *k8s.Cluster
}

func NewServer(cluster *k8s.Cluster) *Server {
	return &Server{
		cluster: cluster,
	}
}

func (s *Server) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}

func (s *Server) EnableOpenTelemetryCollector(ctx context.Context, req *v1.EnableOpenTelemetryCollectorRequest) (*v1.EnableOpenTelemetryCollectorResponse, error) {
	otlp.Enable(ctx, s.cluster)
	return &v1.EnableOpenTelemetryCollectorResponse{}, nil
}

func (s *Server) DisableOpenTelemetryCollector(context.Context, *v1.DisableOpenTelemetryCollectorRequest) (*v1.DisableOpenTelemetryCollectorResponse, error) {
	otlp.Disable(s.cluster)
	return &v1.DisableOpenTelemetryCollectorResponse{}, nil
}

func (s *Server) UpdateOpenTelemetryCollector(context.Context, *v1.UpdateOpenTelemetryCollectorRequest) (*v1.UpdateOpenTelemetryCollectorResponse, error) {
	otlp.Update(s.cluster)
	return &v1.UpdateOpenTelemetryCollectorResponse{}, nil
}
