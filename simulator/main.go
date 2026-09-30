package main

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventsv1 "github.com/bradley-adams/halfback/gen/go/halfback/events/v1"
)

var scoringTypes = []eventsv1.EventType{
	eventsv1.EventType_EVENT_TYPE_TRY,
	eventsv1.EventType_EVENT_TYPE_PENALTY_TRY,
	eventsv1.EventType_EVENT_TYPE_CONVERSION,
	eventsv1.EventType_EVENT_TYPE_PENALTY_GOAL,
	eventsv1.EventType_EVENT_TYPE_DROP_GOAL,
}

func main() {
	matchID := "match-1"
	eventNum := 0

	for {
		eventNum++
		event := &eventsv1.MatchEvent{
			EventId:   fmt.Sprintf("evt-%d", eventNum),
			MatchId:   matchID,
			EventTime: timestamppb.Now(),
			EventType: scoringTypes[rand.Intn(len(scoringTypes))],
		}

		out, err := protojson.Marshal(event)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))

		time.Sleep(2 * time.Second)
	}
}
