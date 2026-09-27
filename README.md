# halfback
Rugby event pipeline: Numaflow → ClickHouse → gRPC → API.

## Local ClickHouse
Container: `halfback-clickhouse` (password: `halfback`)
Start after a reboot: `docker start halfback-clickhouse`
Recreate table from proto: `./scripts/setup-clickhouse.sh`