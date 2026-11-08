package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"

	scoreboardv1connect "github.com/bradley-adams/halfback/gen/go/halfback/scoreboard/v1/scoreboardv1connect"
	"github.com/bradley-adams/halfback/internal/chconn"
)

func main() {
	conn, err := chconn.Open(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := connect.NewServer()
	scoreboardv1connect.RegisterScoreboardServiceHandler(server, NewScoreboardServer(conn))

	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	httpServer := &http.Server{
		Addr:      ":50052",
		Handler:   mux,
		Protocols: protocols,
	}

	log.Println("ScoreboardService listening on :50052")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
