# cdn-load-testing-platform

cdn-load-platform/
├── cmd/
│   ├── controller/
│   │   └── main.go
│   └── agent/
│       └── main.go
│
├── internal/
│   ├── load/
│   │   ├── engine.go
│   │   ├── profiles.go
│   │   ├── token.go
│   │   └── bluegreen.go
│   │
│   ├── metrics/
│   │   └── prometheus.go
│   │
│   └── orchestrator/
│       └── ssm.go
│
├── terraform/
│   ├── controller/
│   │   └── main.tf
│   └── load-nodes/
│       └── main.tf
│
├── docker/
│   ├── agent.Dockerfile
│   └── controller.Dockerfile
│
├── profiles/
│   ├── smoke.json
│   ├── stress.json
│   └── soak.json
│
├── monitoring/
│   ├── prometheus.yml
│   └── grafana-dashboard.json
│
├── go.mod
├── deploy.sh
└── README.md
