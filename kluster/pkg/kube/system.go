package kube

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ultary/monokube/kluster/api/grpc/v1"
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/kube/apps/system/otlp"
)

//var _ v1.SystemServer = System{}

type System struct {
	v1.UnimplementedSystemServer
	cluster *k8s.Cluster
}

func NewSystem(cluster *k8s.Cluster) *System {
	return &System{
		cluster: cluster,
	}
}

func (s *System) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}

func (s *System) EnableOpenTelemetryCollector(ctx context.Context, req *v1.EnableOpenTelemetryCollectorRequest) (*v1.EnableOpenTelemetryCollectorResponse, error) {
	otlp.Enable(ctx, s.cluster)
	return &v1.EnableOpenTelemetryCollectorResponse{}, nil
}

func (s *System) DisableOpenTelemetryCollector(context.Context, *v1.DisableOpenTelemetryCollectorRequest) (*v1.DisableOpenTelemetryCollectorResponse, error) {
	otlp.Disable(s.cluster)
	return &v1.DisableOpenTelemetryCollectorResponse{}, nil
}

func (s *System) UpdateOpenTelemetryCollector(context.Context, *v1.UpdateOpenTelemetryCollectorRequest) (*v1.UpdateOpenTelemetryCollectorResponse, error) {
	otlp.Update(s.cluster)
	return &v1.UpdateOpenTelemetryCollectorResponse{}, nil
}
