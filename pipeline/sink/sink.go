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
)

const insertQuery = `INSERT INTO halfback.match_events SETTINGS format_schema='halfback/events/v1/events.proto:MatchEvent' FORMAT ProtobufSingle`

// ClickHouseSink is a user-defined Numaflow sink that writes MatchEvent
// protobuf bytes straight into ClickHouse via its HTTP interface.
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

// Sink writes each incoming MatchEvent to ClickHouse, one HTTP insert
// per message. Each message is acknowledged or failed independently.
func (s *ClickHouseSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	responses := sinksdk.ResponsesBuilder()

	for datum := range datumStreamCh {
		id := datum.ID()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.insertURL, bytes.NewReader(datum.Value()))
		if err != nil {
			responses = responses.Append(sinksdk.ResponseFailure(id, err.Error()))
			continue
		}

		resp, err := s.client.Do(req)
		if err != nil {
			responses = responses.Append(sinksdk.ResponseFailure(id, err.Error()))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			responses = responses.Append(sinksdk.ResponseFailure(id, string(body)))
			continue
		}
		resp.Body.Close()

		responses = responses.Append(sinksdk.ResponseOK(id))
	}

	return responses
}
