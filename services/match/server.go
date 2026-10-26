package main

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ClickHouse/clickhouse-go/v2"

	matchv1 "github.com/bradley-adams/halfback/gen/go/halfback/match/v1"
)

// MatchServer implements MatchServiceHandler, reading match events
// straight out of ClickHouse.
type MatchServer struct {
	conn clickhouse.Conn
}

func NewMatchServer(conn clickhouse.Conn) *MatchServer {
	return &MatchServer{conn: conn}
}

func (s *MatchServer) ListMatchEvents(
	ctx context.Context,
	req *connect.Request[matchv1.ListMatchEventsRequest],
) (*connect.Response[matchv1.ListMatchEventsResponse], error) {
	rows, err := s.conn.Query(ctx, `
		SELECT event_id, event_type, event_time.seconds
		FROM halfback.match_events FINAL
		WHERE match_id = $1
		ORDER BY event_time.seconds
	`, req.Msg.MatchId)
	if err != nil {
		return nil, fmt.Errorf("query match events: %w", err)
	}
	defer rows.Close()

	resp := &matchv1.ListMatchEventsResponse{}

	for rows.Next() {
		var (
			eventID   string
			eventType string
			seconds   int64
		)
		if err := rows.Scan(&eventID, &eventType, &seconds); err != nil {
			return nil, fmt.Errorf("scan match event row: %w", err)
		}

		resp.Events = append(resp.Events, &matchv1.MatchEventSummary{
			EventId:          eventID,
			EventType:        eventType,
			EventTimeSeconds: seconds,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate match event rows: %w", err)
	}

	return connect.NewResponse(resp), nil
}
