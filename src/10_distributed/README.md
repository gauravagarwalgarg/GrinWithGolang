# Module 10: Distributed Systems in Go

> Real-world infrastructure patterns: gRPC, messaging, databases, orchestration, and observability.

## Submodules

| Folder | Tech | What You Build |
|--------|------|----------------|
| `grpc_service/` | gRPC + Protobuf | User service with unary + streaming RPCs |
| `messaging/nsq/` | NSQ | Pub/sub event system with producer + consumer |
| `messaging/kafka/` | Kafka | Event-driven order pipeline with partitions |
| `databases/postgres/` | PostgreSQL | Repository pattern with sqlx, migrations |
| `databases/mongo/` | MongoDB | Aggregation pipelines, change streams |
| `observability/` | Prometheus + Tempo | Metrics, distributed tracing, Grafana dashboards |
| `kubernetes/` | K8s manifests | Deployments, services, HPA, ConfigMaps |

## Architecture Overview

```
┌─────────────┐     gRPC      ┌─────────────┐
│  API Gateway├──────────────►│ User Service │
└──────┬──────┘               └──────┬──────┘
       │                             │
       │ HTTP                        │ DB
       ▼                             ▼
┌─────────────┐              ┌─────────────┐
│   Client    │              │  PostgreSQL  │
└─────────────┘              └─────────────┘
       │
       │ Events
       ▼
┌─────────────┐   consume    ┌─────────────┐
│    Kafka    ├─────────────►│  Worker Svc  │
└─────────────┘              └──────┬──────┘
                                    │
                                    │ trace
                                    ▼
                             ┌─────────────┐
                             │ Tempo/Prom  │
                             └─────────────┘
```

## Quick Start

Each submodule is self-contained with its own `docker-compose.yaml` and `Makefile`.

```bash
# gRPC service
cd grpc_service && make proto && make run

# Kafka pipeline
cd messaging/kafka && docker-compose up -d && make run

# Full observability stack
cd observability && docker-compose up -d
```
