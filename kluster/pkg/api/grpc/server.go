package grpc

import (
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/kube"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/ultary/monokube/kluster/api/grpc/v1"
)

type Server struct {
	v1.KlusterServer
	v1.SystemServer

	cluster *kube.Cluster
	server  *grpc.Server
}

func NewServer(client *k8s.Client) *Server {

	server := grpc.NewServer()
	retval := &Server{
		cluster: kube.NewCluster(client),
		server:  server,
	}

	v1.RegisterKlusterServer(server, retval)
	v1.RegisterSystemServer(server, retval)
	return retval
}

func (s *Server) Serve(network, address string) error {
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Errorf("[gRPC] Failed to listen: %v", err)
		return err
	}
	if err = s.server.Serve(lis); err != nil {
		log.Errorf("[gRPC] Failed to serve: %v", err)
		return err
	}
	return nil
}

func (s *Server) Stop() {
	log.Info("[gRPC] Stopping server")
	s.server.Stop()
	log.Info("[gRPC] Stopped server")
}
