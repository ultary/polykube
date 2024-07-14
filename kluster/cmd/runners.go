package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ultary/monokube/kluster/pkg/api/grpc"
	"github.com/ultary/monokube/kluster/pkg/api/http"
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/kube/watch"
)

type Runner interface {
	Use() string
	Short() string
	Run(cmd *cobra.Command, args []string)
}

////////////////////////////////////////////////////////////////
//
//  Installer
//

type installer struct {
}

func (i *installer) Use() string {
	return "install"
}

func (i *installer) Short() string {
	return "Run monokube's kluster installer"
}

func (i *installer) Run(cmd *cobra.Command, args []string) {
	log.Info("Install")
}

////////////////////////////////////////////////////////////////
//
//  Server
//

type server struct {
	incluster   bool
	kubeconfig  string
	kubecontext string
}

func (s *server) Use() string {
	return "serve"
}

func (s *server) Short() string {
	return "Run monokube's kluster server"
}

func (s *server) Run(cmd *cobra.Command, args []string) {

	////////////////////////////////////////////////////////////////
	//
	//  PostgreSQL
	//

	// https://www.npgsql.org/doc/connection-string-parameters.html
	const pgdsn = "postgres://postgres:postgrespw@localhost:5432/postgres?application_name=ultary&sslmode=disable&connect_timeout=5"
	pgconfig, err := pgxpool.ParseConfig(pgdsn)
	if err != nil {
		log.Fatal(err)
	}
	pgpool, err := pgxpool.NewWithConfig(context.Background(), pgconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pgpool.Close()

	////////////////////////////////////////////////////////////////
	//
	//  Server Instances
	//

	client := k8s.NewClient(s.kubeconfig, s.kubecontext)
	grpcServer := grpc.NewServer(client)
	httpServer := http.NewServer()
	watchTower := watch.NewTower(client)

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		_ = <-sigs

		log.Info("Gracefully shutting down ...")
		grpcServer.Stop()
		httpServer.Shutdown()
		watchTower.Shutdown()
	}()

	// WatchTower
	wg.Add(1)
	go func() {
		defer wg.Done()

		watchTower.Watch()
	}()

	// HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()

		address := "0.0.0.0:9090"
		if err := httpServer.Listen(address); err != nil {
			log.Fatalf("Failed to listen http server: %v", err)
		}
	}()

	// gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()

		network, address := "tcp4", "0.0.0.0:50051"
		//network, address := "unix", "/tmp/kluster.sock"
		if err := grpcServer.Serve(network, address); err != nil {
			log.Fatalf("Failed to serve grpc server: %v", err)
		}
	}()

	wg.Wait()
}
