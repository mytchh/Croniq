# Croniq
Web UI for managing k8s cron jobs
croniq/
├── frontend/
│   ├── package.json
│   ├── src/
│   │   ├── App.tsx
│   │   └── components/
│   │       └── CronJobList.tsx
├── backend/
│   ├── main.go
│   ├── handlers/
│   │   └── cronjob.go
│   └── k8s/
│       └── client.go
└── go.mod