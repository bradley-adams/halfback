package main

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventsv1 "github.com/bradley-adams/halfback/gen/go/halfback/events/v1"
)

func main() {
	event := &eventsv1.MatchEvent{
		EventId:   "evt-1",
		MatchId:   "match-1",
		EventTime: timestamppb.Now(),
		EventType: eventsv1.EventType_EVENT_TYPE_TRY,
	}

	data, err := proto.Marshal(event)
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal error:", err)
		os.Exit(1)
	}

	if err := os.WriteFile("event.bin", data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}

	fmt.Println("wrote event.bin,", len(data), "bytes")
}
