# halfback
Rugby event pipeline: Numaflow → ClickHouse → gRPC services → API (in progress).

## Local ClickHouse
Container: `halfback-clickhouse` (password: `halfback`)
Start after a reboot: `docker start halfback-clickhouse`
Recreate table from proto: `./scripts/setup-clickhouse.sh`

## Local Kubernetes (kind) + Numaflow
Cluster: `halfback` (created via `deploy/kind/kind-config.yaml`)
Recreate cluster: `kind create cluster --config deploy/kind/kind-config.yaml`
Install Numaflow into the cluster: `./scripts/setup-numaflow.sh`

### The pipeline
Manifest: `deploy/numaflow/halfback-pipeline.yaml` (source: simulator, sink: ClickHouse writer)
Apply/update: `kubectl apply -f deploy/numaflow/halfback-pipeline.yaml`
Check status: `kubectl get pods`

Rebuild and reload an image after a code change, then restart its pod:
```bash
docker build -f pipeline/source/Dockerfile -t halfback-simulator:dev .   # or pipeline/sink, halfback-sink
kind load docker-image halfback-simulator:dev --name halfback            # match the tag above
kubectl delete pod -l numaflow.numaproj.io/vertex-name=in                # or vertex-name=out for the sink
```

ClickHouse is reached from inside the cluster via a Service/Endpoints pointing at the Docker `kind` network's gateway IP (see `deploy/clickhouse/service.yaml`). If the kind cluster is ever recreated, that IP may need updating — check with:
`docker network inspect kind -f '{{json .IPAM.Config}}'`

## MatchService (Connect — gRPC, gRPC-Web, and plain HTTP/JSON)
Run locally: `go run ./services/match` (listens on :50051)
Requires ClickHouse running (see above).

Test with grpcurl (works unmodified against Connect):
grpcurl -plaintext -import-path proto -proto halfback/match/v1/match.proto \
  -d '{"match_id": "match-1"}' \
  localhost:50051 halfback.match.v1.MatchService/ListMatchEvents

Or test with plain curl, no gRPC tooling required:
curl -s -X POST http://localhost:50051/halfback.match.v1.MatchService/ListMatchEvents \
  -H "Content-Type: application/json" \
  -d '{"match_id": "match-1"}'