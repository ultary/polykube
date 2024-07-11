package main

import (
	"context"
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/k8s/keeper"
	"log"
	"net"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ultary/monokube/kluster/api/grpc/v1"
)

func main() {

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

	var wg sync.WaitGroup

	ctx := k8s.NewContext(context.Background())
	k := keeper.Init(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()
		k.Listen()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := grpc.NewServer()
		v1.RegisterSystemServiceServer(server, &System{})

		lis, err := net.Listen("tcp4", "127.0.0.1:50051")
		//lis, err := net.Listen("unix", "/tmp/kluster.sock")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err = server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}

type System struct {
	v1.SystemServiceServer
}

func (s *System) Ping(ctx context.Context, empty *emptypb.Empty) (*v1.Pong, error) {
	return &v1.Pong{
		Pong: "pong",
	}, nil
}
