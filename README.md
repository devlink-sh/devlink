# devlink

## Folder Structure

```plaintext
devlink/
├── cmd/
│   └── devlink/
│       ├── git/
│       │   ├── connect.go
│       │   ├── git.go
│       │   ├── serve.go
│       │   └── utils.go
│       │
│       ├── main.go
│       └── root.go
│
├── internal/
│   └── git/
│       ├── process/
│       │   └── manager.go
│       │
│       └── proxy/
│           └── tcp.go
│
├── .gitignore
├── go.mod
├── go.sum
└── README.md




A minimal implementation of `devlink git serve` and `devlink git connect` that
wraps a local `git daemon` behind an OpenZiti secure service.


> **Important:** Do not commit your OpenZiti identity JSON files into this repo. Keep them private.


## Quick build


```bash
go mod tidy
go build -o devlink ./cmd/devlink
```








## Serve (on machine A)


```bash
# inside a git repo
./devlink git serve --identity /path/to/peerA.json
 - validates repo and identity
  - starts local git daemon on ephemeral localhost port
  - binds an OpenZiti service and proxies incoming Ziti connections to git daemon
  - process-group-aware graceful shutdown (SIGTERM -> SIGKILL escalation)

```


## Connect (on machine B)


```bash
# inside a git repo
./devlink git connect <service-name> --identity /path/to/peerB.json --name ananaya-wip
# then use git fetch/pull/push with the configured remote
- validates repo and identity
  - validates connectivity to service via OpenZiti
  - opens a loopback listener and proxies local git client connections to the service
  - configures git remote automatically (git remote add/set-url)

```
***prototype to demonstrate the p2p git remote concept using openziti
