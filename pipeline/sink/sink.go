package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"

	eventsv1 "github.com/bradley-adams/halfback/gen/go/halfback/events/v1"
)

const insertQuery = `INSERT INTO halfback.match_events SETTINGS format_schema='halfback/events/v1/events.proto:MatchEvent' FORMAT Protobuf`

type ClickHouseSink struct {
	insertURL string
	client    *http.Client
}

func NewClickHouseSink() *ClickHouseSink {
	host := getEnv("CLICKHOUSE_HOST", "localhost:8123")
	password := getEnv("CLICKHOUSE_PASSWORD", "halfback")

	params := url.Values{}
	params.Set("user", "default")
	params.Set("password", password)
	params.Set("query", insertQuery)

	insertURL := fmt.Sprintf("http://%s/?%s", host, params.Encode())

	return &ClickHouseSink{
		insertURL: insertURL,
		client:    &http.Client{},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Batches all messages into one length-delimited Protobuf stream and
// sends a single insert. The HTTP request either succeeds or fails as
// a whole; individually malformed messages are excluded before that.
func (s *ClickHouseSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	responses := sinksdk.ResponsesBuilder()

	var buf bytes.Buffer
	var ids []string

	for datum := range datumStreamCh {
		event := &eventsv1.MatchEvent{}
		if err := proto.Unmarshal(datum.Value(), event); err != nil {
			responses = responses.Append(sinksdk.ResponseFailure(datum.ID(), err.Error()))
			continue
		}

		if _, err := protodelim.MarshalTo(&buf, event); err != nil {
			responses = responses.Append(sinksdk.ResponseFailure(datum.ID(), err.Error()))
			continue
		}

		ids = append(ids, datum.ID())
	}

	if len(ids) == 0 {
		return responses
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.insertURL, &buf)
	if err != nil {
		for _, id := range ids {
			responses = responses.Append(sinksdk.ResponseFailure(id, err.Error()))
		}
		return responses
	}

	resp, err := s.client.Do(req)
	if err != nil {
		for _, id := range ids {
			responses = responses.Append(sinksdk.ResponseFailure(id, err.Error()))
		}
		return responses
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		for _, id := range ids {
			responses = responses.Append(sinksdk.ResponseFailure(id, string(body)))
		}
		return responses
	}

	for _, id := range ids {
		responses = responses.Append(sinksdk.ResponseOK(id))
	}

	return responses
}
