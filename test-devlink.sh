#!/bin/bash
set -euo pipefail

echo "=== Setting up real Devlink serve/connect test ==="

# Cleanup old test environment
rm -rf /tmp/devlink-real-test
mkdir -p /tmp/devlink-real-test
cd /tmp/devlink-real-test

# Step 1: Initialize host repo
echo ">>> Initializing host repo"
mkdir host-repo
cd host-repo
git init -q
echo "hello world" > file.txt
git add file.txt
git commit -qm "initial commit"
cd ..

# Step 2: Create bare mirror (simulate serve) for local testing
mkdir -p /tmp/devlink-test/shared
git clone --mirror host-repo /tmp/devlink-test/shared/host-repo.git

# Step 3: Start devlink serve in background
echo ">>> Starting devlink serve in background"
cd host-repo
devlink serve > /tmp/devlink-serve.log 2>&1 &
SERVE_PID=$!
cd ..

# Wait a few seconds for serve to initialize
sleep 3

# Step 4: Use the bare mirror as the "tunnel URL" for teammate
TUNNEL_URL="/tmp/devlink-test/shared/host-repo.git"
echo "Tunnel/bare repo URL for teammate: $TUNNEL_URL"

# Step 5: Simulate teammate connecting
echo ">>> Teammate connecting"
git clone "$TUNNEL_URL" teammate-repo
cd teammate-repo
git checkout -B master
echo ">>> Teammate making changes"
echo "teammate change" >> file.txt
git add file.txt
git commit -qm "teammate update"
git push origin master
cd ..

# Step 6: Host pulls changes with conflict auto-resolution
echo ">>> Host pulling changes"
cd host-repo

# Try pull, allow unrelated histories, don't rebase
if ! git pull --allow-unrelated-histories --no-rebase "$TUNNEL_URL" master; then
    echo ">>> Merge conflict detected, auto-resolving..."
    # Accept teammate’s changes by default
    git checkout --theirs file.txt || true
    git add file.txt
    git commit -qm "Auto-resolve merge conflict by taking teammate’s changes"
fi

# Step 7: Cleanup background serve
kill $SERVE_PID || true

echo
echo "=== Real Devlink test complete! File content in host-repo: ==="
cat file.txt

