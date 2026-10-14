# halfback
Rugby event pipeline: Numaflow → ClickHouse → gRPC → API.

## Local ClickHouse
Container: `halfback-clickhouse` (password: `halfback`)
Start after a reboot: `docker start halfback-clickhouse`
Recreate table from proto: `./scripts/setup-clickhouse.sh`

## MatchService (gRPC)
Run locally: `go run ./services/match` (listens on :50051)
Requires ClickHouse running (see above).

Test with grpcurl:
grpcurl -plaintext -import-path proto -proto halfback/match/v1/match.proto \
  -d '{"match_id": "match-1"}' \
  localhost:50051 halfback.match.v1.MatchService/ListMatchEvents