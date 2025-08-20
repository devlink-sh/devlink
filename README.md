# devlink

## Folder Structure

devlink/
├── cmd/
│   ├── root.go
│   ├── env.go
│   ├── git.go
│   └── db.go
│
├── internal/
│   ├── p2p/
│   │   ├── client.go
│   │   └── server.go
│   │
│   ├── git/
│   │   └── daemon.go
│   │
│   └── util/
│       └── token.go
│
├──.gitignore
├── go.mod
├── go.sum
└── main.go