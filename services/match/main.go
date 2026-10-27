package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	mux := http.NewServeMux()
	path, handler := matchv1connect.NewMatchServiceHandler(NewMatchServer(conn))
	mux.Handle(path, handler)

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	server := &http.Server{
		Addr:      ":50051",
		Handler:   mux,
		Protocols: protocols,
	}

	log.Println("MatchService listening on :50051")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
