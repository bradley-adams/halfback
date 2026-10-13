package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"google.golang.org/grpc"

	matchv1 "github.com/bradley-adams/halfback/gen/go/halfback/match/v1"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	host := getEnv("CLICKHOUSE_HOST", "localhost:9000")
	password := getEnv("CLICKHOUSE_PASSWORD", "halfback")

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: "halfback",
			Username: "default",
			Password: password,
		},
	})
	if err != nil {
		log.Fatal("failed to connect to clickhouse: ", err)
	}
	defer conn.Close()

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping clickhouse: ", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	grpcServer := grpc.NewServer()
	matchv1.RegisterMatchServiceServer(grpcServer, NewMatchServer(conn))

	log.Println("MatchService listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
