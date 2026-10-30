package main

import (
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/vanguard"

	apiv1connect "github.com/bradley-adams/halfback/gen/go/halfback/api/v1/apiv1connect"
	matchv1connect "github.com/bradley-adams/halfback/gen/go/halfback/match/v1/matchv1connect"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	matchServiceURL := getEnv("MATCH_SERVICE_URL", "http://localhost:50051")

	matchClient := matchv1connect.NewMatchServiceClient(
		connect.NewClient(connecthttp.NewTransport(http.DefaultClient, matchServiceURL)),
	)

	publicServer := NewMatchPublicServer(matchClient)

	server := connect.NewServer()
	apiv1connect.RegisterMatchPublicServiceHandler(server, publicServer)

	mux := http.NewServeMux()
	connecthttp.Mount(mux, server) // Connect, gRPC, gRPC-Web
	if err := vanguard.Mount(mux, server); err != nil {
		log.Fatal("failed to mount vanguard: ", err)
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	httpServer := &http.Server{
		Addr:      ":8080",
		Handler:   mux,
		Protocols: protocols,
	}

	log.Println("API listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
