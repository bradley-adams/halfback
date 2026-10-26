package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

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

	log.Println("MatchService listening on :50051")
	err = http.ListenAndServe(
		":50051",
		// h2c serves HTTP/2 without TLS, so plain gRPC clients still work locally.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
