package main

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"

	scoreboardv1 "github.com/bradley-adams/halfback/gen/go/halfback/scoreboard/v1"
)

// ScoreboardServer implements ScoreboardServiceHandler, computing
// per-match totals straight out of ClickHouse.
type ScoreboardServer struct {
	conn clickhouse.Conn
}

func NewScoreboardServer(conn clickhouse.Conn) *ScoreboardServer {
	return &ScoreboardServer{conn: conn}
}

// Points per event type live here, not in the table: try 5, conversion 2,
// penalty goal 3, drop goal 3, penalty try 7. The toInt64 casts matter
// because count() and sum() return unsigned types, and the driver scans
// integer types strictly.
const listMatchesQuery = `
	SELECT
		match_id,
		toInt64(count()) AS event_count,
		toInt64(sum(multiIf(
			event_type = 'EVENT_TYPE_TRY', 5,
			event_type = 'EVENT_TYPE_CONVERSION', 2,
			event_type = 'EVENT_TYPE_PENALTY_GOAL', 3,
			event_type = 'EVENT_TYPE_DROP_GOAL', 3,
			event_type = 'EVENT_TYPE_PENALTY_TRY', 7,
			0
		))) AS total_points
	FROM halfback.match_events FINAL
	GROUP BY match_id
	ORDER BY match_id
`

func (s *ScoreboardServer) ListMatches(
	ctx context.Context,
	_ *scoreboardv1.ListMatchesRequest,
) (*scoreboardv1.ListMatchesResponse, error) {
	rows, err := s.conn.Query(ctx, listMatchesQuery)
	if err != nil {
		return nil, fmt.Errorf("query matches: %w", err)
	}
	defer rows.Close()

	resp := &scoreboardv1.ListMatchesResponse{}

	for rows.Next() {
		var (
			matchID     string
			eventCount  int64
			totalPoints int64
		)
		if err := rows.Scan(&matchID, &eventCount, &totalPoints); err != nil {
			return nil, fmt.Errorf("scan match row: %w", err)
		}

		resp.Matches = append(resp.Matches, &scoreboardv1.MatchScore{
			MatchId:     matchID,
			EventCount:  eventCount,
			TotalPoints: totalPoints,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate match rows: %w", err)
	}

	return resp, nil
}
