package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"

	matchv1connect "github.com/bradley-adams/halfback/gen/go/halfback/match/v1/matchv1connect"
	"github.com/bradley-adams/halfback/internal/chconn"
)

func main() {
	conn, err := chconn.Open(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

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
