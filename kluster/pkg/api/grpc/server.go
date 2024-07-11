package grpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/ultary/monokube/kluster/pkg/api/grpc/v1"
)

type Server struct {
	server *grpc.Server
}

func NewServer() *Server {

	server := grpc.NewServer()
	v1.RegisterKlusterServer(server)

	return &Server{
		server: server,
	}
}

func (s *Server) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}
