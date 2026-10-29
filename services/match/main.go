package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/ClickHouse/clickhouse-go/v2"

	matchv1connect "github.com/bradley-adams/halfback/gen/go/halfback/match/v1/matchv1connect"
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

	server := connect.NewServer()
	matchv1connect.RegisterMatchServiceHandler(server, NewMatchServer(conn))

	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	httpServer := &http.Server{
		Addr:      ":50051",
		Handler:   mux,
		Protocols: protocols,
	}

	log.Println("MatchService listening on :50051")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
