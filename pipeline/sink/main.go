package main

import (
	"context"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

func main() {
	err := sinksdk.NewServer(NewClickHouseSink()).Start(context.Background())
	if err != nil {
		log.Panic("failed to start sink server: ", err)
	}
}
