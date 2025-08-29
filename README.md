# devlink

## Folder Structure

```plaintext
devlink/
├── cmd/
│   └── root.go         # entrypoint: parses top-level args, dispatches to subcommands
│   └── serve.go        # implementation of `devlink-git serve`
│   └── connect.go      # implementation of `devlink-git connect`
│
├── internal/
│   ├── gitutils/
│   │   └── git.go      # helpers: repo detection, runStreaming, runSilently
│   ├── zrok/
│   │   └── tunnel.go   # helpers: start zrok, stream hints, extract ports
│   └── signal/
│       └── interrupt.go # helper: waitForInterrupt
│
├── go.mod              # module definition (e.g. module github.com/you/devlink-git)
├── go.sum
└── main.go 
└── README.md



// devlink git — minimal MVP that wraps `git daemon` behind a zrok tunnel
//
// Commands:
// devlink-git serve # run inside a git repo (creates a bare mirror, serves it, exposes via zrok)
// devlink-git connect --token <token> # access the share and print clone/push instructions
//
// Prereqs on both machines:
// - git installed
// - zrok installed & "zrok enable" done once (account/activation)
//
// SECURITY NOTE: This MVP is for ephemeral, trusted use during crunch time.
// It enables receive-pack (push) on the served bare repo. Use zrok private shares and rotate tokens.


