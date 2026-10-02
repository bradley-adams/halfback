package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/sourcer"
	"google.golang.org/protobuf/proto"
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

// MatchEventSource is a user-defined Numaflow source that generates
// fake MatchEvent messages. It has no real backlog or durability, since
// every message is invented on the fly.
type MatchEventSource struct {
	matchID  string
	eventNum int
}

func NewMatchEventSource() *MatchEventSource {
	return &MatchEventSource{matchID: "match-1"}
}

// Read generates up to readRequest.Count() messages and sends them to messageCh.
// It must return once that many messages are sent, or the timeout elapses.
func (s *MatchEventSource) Read(ctx context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- sourcesdk.Message) {
	deadline := time.Now().Add(readRequest.TimeOut())

	for i := 0; i < int(readRequest.Count()); i++ {
		if time.Now().After(deadline) {
			return
		}

		s.eventNum++
		event := &eventsv1.MatchEvent{
			EventId:   fmt.Sprintf("evt-%d", s.eventNum),
			MatchId:   s.matchID,
			EventTime: timestamppb.Now(),
			EventType: scoringTypes[rand.Intn(len(scoringTypes))],
		}

		data, err := proto.Marshal(event)
		if err != nil {
			continue
		}

		messageCh <- sourcesdk.NewMessage(
			data,
			sourcesdk.NewOffset([]byte(event.EventId), 0),
			event.EventTime.AsTime(),
		)

		time.Sleep(2 * time.Second)
	}
}

// Ack acknowledges messages the pipeline has finished processing. Since
// generated messages aren't stored anywhere, there's nothing to do.
func (s *MatchEventSource) Ack(_ context.Context, _ sourcesdk.AckRequest) {}

// Nack negatively acknowledges offsets the pipeline failed to process.
// Since generated messages aren't stored anywhere, there's nothing to do.
func (s *MatchEventSource) Nack(_ context.Context, _ sourcesdk.NackRequest) {}

// Pending reports how many messages are waiting to be read. We generate
// on demand, so there's never a backlog.
func (s *MatchEventSource) Pending(_ context.Context) int64 {
	return 0
}

// ActivePartitions returns the partitions this source currently reads
// from. We only have one logical stream, so we use the SDK's default.
func (s *MatchEventSource) ActivePartitions(_ context.Context) []int32 {
	return sourcesdk.DefaultPartitions()
}

// TotalPartitions returns the total number of partitions this source
// has. We don't have a meaningful concept of partitions, so we report
// nothing, as the interface allows.
func (s *MatchEventSource) TotalPartitions(_ context.Context) *int32 {
	return nil
}
