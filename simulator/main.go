package main

import (
	"context"
	"log"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/sourcer"
)

func main() {
	err := sourcesdk.NewServer(NewMatchEventSource()).Start(context.Background())
	if err != nil {
		log.Panic("failed to start source server: ", err)
	}
}
