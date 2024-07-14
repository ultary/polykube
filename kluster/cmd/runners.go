package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ultary/monokube/kluster/pkg/kube"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// ---- postgresql ----

	// https://www.npgsql.org/doc/connection-string-parameters.html
	const pgdsn = "postgres://postgres:postgrespw@localhost:5432/postgres?application_name=ultary&sslmode=disable&connect_timeout=5"
	//dsn := "user=youruser password=yourpassword dbname=yourdb host=localhost port=5432 sslmode=disable"
	pgxconfig, err := pgxpool.ParseConfig(pgdsn)
	if err != nil {
		log.Fatal(err)
	}
	pgxpool, err := pgxpool.NewWithConfig(context.Background(), pgxconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxpool.Close()

	// --- gorm ----

	// logger := logger.New(
	// 	logrus.NewWriter(),
	// 	logger.Config{
	// 		SlowThreshold: time.Millisecond,
	// 		LogLevel:      logger.Warn,
	// 		Colorful:      false,
	// 	},
	// )

	dsn := "host=localhost user=postgres password=postgrespw dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Seoul"
	dialector := postgres.Open(dsn)
	gormConfig := &gorm.Config{
		// Logger: logger,
	}
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatal(err)
	}
	// if err = db.Use(tracing.NewPlugin()); err != nil {
	// 	log.Fatal(err)
	// }

	////////////////////////////////////////////////////////////////
	//
	//  Server Instances
	//

	client := k8s.NewClient(s.kubeconfig, s.kubecontext)
	cluster := k8s.NewCluster(client, db)
	system := kube.NewSystem(cluster)

	grpcServer := grpc.NewServer()
	grpcServer.RegisterSystemServer(system)

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
