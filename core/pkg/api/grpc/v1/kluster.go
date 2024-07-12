package v1

import (
	"context"

	"google.golang.org/grpc"

	v1 "github.com/ultary/monokube/core/api/grpc/v1"
)

func RegisterKlusterServer(s grpc.ServiceRegistrar) {
	v1.RegisterKlusterServer(s, &Kluster{})
}

type Kluster struct {
	v1.KlusterServer
}

func (s *Kluster) SyncOpenTelemetry(context.Context, *v1.SyncOpenTelemetryRequest) (*v1.SyncOpenTelemetryResponse, error) {
	return nil, nil
}
