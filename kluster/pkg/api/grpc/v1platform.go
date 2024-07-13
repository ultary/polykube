package grpc

import (
	"context"

	v1 "github.com/ultary/monokube/kluster/api/grpc/v1"
)

func (s *Server) CreateNamespace(ctx context.Context, req *v1.CreateNamespaceRequest) (*v1.CreateNamespaceResponse, error) {
	return &v1.CreateNamespaceResponse{}, nil
}
