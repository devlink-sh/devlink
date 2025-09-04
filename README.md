
Secure, peer-to-peer developer collaboration.
Share environments, databases, repos, and artifacts without pushing to the cloud.
# **DevLink CLI Reference**

> A comprehensive guide to all DevLink commands for secure, peer-to-peer developer collaboration.

---

## **📦 Installation**

Build and install DevLink locally:

```bash
git clone <devlink-repo-url>
cd DevLink
go build -o devlink
```

Optionally, move it to PATH for global usage:

```bash
sudo mv devlink /usr/local/bin/
```

---

## **🔐 Environment Sharing (`env`)**

> Share `.env` or secret files securely with teammates.

**Share an environment:**

```bash
devlink env share
```

**Fetch an environment:**

```bash
devlink env get <share-token>
```

**Key Notes:**

* ✅ End-to-end encrypted
* ✅ Instant transfer
* ✅ No third-party servers

---

## **🗄️ Database Sharing (`db`)**

> Grant ephemeral, read-only access to local databases.

**Share a database (example port 5432):**

```bash
devlink db share 5432
```

**Connect to a shared DB locally:**

```bash
devlink db get <share-token> <local-port>
```

**Example (Postgres):**

```bash
psql -h 127.0.0.1 -p <local-port> -U <db-username> -d <db-name>
```

**Benefits:**

* ✅ Live integration
* ✅ Secure and ephemeral
* ✅ Minimal setup

---

## **🌱 Git Repository Sharing (`git`)**

> Share your local Git repository without pushing WIP code.

**Serve a repo:**

```bash
cd ~/projects/my-repo
devlink git serve .
```

**Connect to a shared repo:**

```bash
devlink git connect <share-token> my-repo.git
```

**Clone a shared repo:**

```bash
git clone git://127.0.0.1:9418/my-repo.git
```

**Continue normal workflow:**

```bash
git add <file>
git commit -m "message"
git push origin main
git pull
```

---

## **🧩 Directory Sharing (`dir`)** *

> Share entire local directories with teammates in a peer-to-peer manner.

**Share a directory:**

```bash
devlink dir share <path>
```

**Connect / fetch a shared directory:**

```bash
devlink dir get <share-token> <local-path>
```

---

## **🚀 Hive: Ephemeral Staging Environments (`hive`)**

> Spin up shared temporary environments for integrated testing.

**Create a new Hive:**

```bash
devlink hive create <hive-name>
```

**Contribute a service:**

```bash
devlink hive contribute --service <name> --port <port> --hive <invite-token>
```

**Connect to an existing Hive:**

```bash
devlink hive connect --hive <invite-token>
```

**Teardown:**

* Press `Ctrl+C` to stop contributing
* Environment disappears automatically

**Impact:**

* ✅ Live, ephemeral multi-service environment
* ✅ Fast testing & debugging
* ✅ Zero-trust P2P network

---

## **🔗 Pairing with Hive Controller (`pair`)**

> Connect your CLI to the central Hive Controller for coordination.

**Pair CLI with controller:**

```bash
devlink pair get <pair-token> <controller-port>
```

* Controller becomes available at: `http://localhost:<controller-port>`
* Required before creating or connecting to Hives

---

## **📦 Registry Sharing (`registry`)**

> Share Docker images or local artifacts P2P with teammates.

**Send a Docker image / artifact:**

```bash
devlink registry send <image-name>
```

**Receive an image / artifact:**

```bash
devlink registry receive <share-token>
```

**Key Notes:**

* ✅ Direct P2P transfer
* ✅ Faster than central registries
* ✅ Works offline / within team network

---

## **📊 Why DevLink is Different**

* 🔒 Zero-trust: each share scoped to a single tunnel
* ⚡ Peer-to-peer: no cloud hosting required
* 🔐 Secure by design: `.env` files, DBs, Git branches never leak
* 🎯 Real-world usage: hackathons, live demos, team sprints

---

## **💡 Example Workflow**

1. Pair CLI: `devlink pair get <token> <port>`
2. Share environment: `devlink env share` → teammate `devlink env get <token>`
3. Share DB: `devlink db share <port>` → teammate connects
4. Share Git repo: `devlink git serve` → teammate clones
5. Spin up ephemeral staging: `devlink hive create` → contribute services → QA connects
6. Share artifacts/images: `devlink registry send <image>` → teammate receives
7. Teardown: Ctrl+C → environment disappears

---

✅ **Everything is generic & reusable:**

* Tokens, ports, paths, usernames are placeholders
* No personal info is exposed

---

