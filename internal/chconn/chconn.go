package chconn

import (
	"context"
	"fmt"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Open connects to ClickHouse using CLICKHOUSE_HOST and CLICKHOUSE_PASSWORD
// and pings it, so a service fails at startup rather than on its first query.
func Open(ctx context.Context) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{getEnv("CLICKHOUSE_HOST", "localhost:9000")},
		Auth: clickhouse.Auth{
			Database: "halfback",
			Username: "default",
			Password: getEnv("CLICKHOUSE_PASSWORD", "halfback"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open clickhouse: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping clickhouse: %w", err)
	}

	return conn, nil
}
