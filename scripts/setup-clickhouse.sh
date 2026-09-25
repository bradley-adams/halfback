#!/usr/bin/env bash
set -euo pipefail

CONTAINER=halfback-clickhouse
SCHEMA_DIR=/var/lib/clickhouse/format_schemas

echo "Copying proto files into ClickHouse container..."
docker exec -it "$CONTAINER" rm -rf "$SCHEMA_DIR/halfback"
docker exec -it "$CONTAINER" mkdir -p "$SCHEMA_DIR/halfback/events/v1"
docker cp proto/halfback/events/v1/events.proto \
  "$CONTAINER:$SCHEMA_DIR/halfback/events/v1/events.proto"

echo "Linking well-known protobuf types..."
docker exec -it "$CONTAINER" sh -c \
  "[ -L $SCHEMA_DIR/google ] || ln -s /usr/share/clickhouse/protos/google $SCHEMA_DIR/google"

echo "Ensuring placeholder file for schema-only queries exists..."
docker exec -it "$CONTAINER" touch /var/lib/clickhouse/user_files/nonexist

echo "Creating database and table..."
docker exec -it "$CONTAINER" clickhouse-client --multiquery --query "
CREATE DATABASE IF NOT EXISTS halfback;

CREATE TABLE IF NOT EXISTS halfback.match_events
ENGINE = MergeTree
ORDER BY (match_id, event_id)
AS SELECT * FROM file('nonexist', 'Protobuf')
SETTINGS format_schema = 'halfback/events/v1/events.proto:MatchEvent';
"

echo "Done. Current structure:"
docker exec -it "$CONTAINER" clickhouse-client --query "DESCRIBE halfback.match_events" --format PrettyCompact