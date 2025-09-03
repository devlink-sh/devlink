
# DevLink

DevLink is a CLI tool that enables seamless, **peer-to-peer sharing of local Git repositories** without exposing your machine to the public internet. This guide covers installation, repository sharing, and normal workflow.

---

## **Folder Structure**

```plaintext
DevLink/
├─ cmd/
│  ├─ db/
│  │  ├─ db.go
│  │  ├─ get.go
│  │  └─ share.go
│  ├─ env/
│  │  ├─ env.go
│  │  ├─ get.go
│  │  └─ share.go
│  ├─ git/
│  │  ├─ connect.go       
│  │  ├─ git.go          
│  │  └─ serve.go        
│  ├─ pair/
│  │  ├─ get.go
│  │  ├─ pair.go
│  │  └─ share.go
│  └─ registry/
│     ├─ get.go
│     ├─ registry.go
│     └─ share.go
├─ internal/
│  └─ proxy.go            
├─ main.go
├─ go.mod
├─ go.sum
└─ README.md
````


---

# DevLink

DevLink is a CLI tool that enables **seamless, peer-to-peer sharing of local Git repositories** without exposing your machine to the public internet.

With DevLink you can turn *any local Git repo* into a temporary share, and your teammates can instantly clone or push without needing GitHub, VPNs, or public servers.

---

## 📦 Installation

1. **Clone the DevLink repository**

```bash
git clone <devlink-repo-url>
cd DevLink
```

2. **Build the CLI tool**

```bash
go build -o devlink
```

This generates a `devlink` executable in your current directory.

---

## 🚀 Usage

### 1. Start sharing a repository (Host machine)

From inside the repo you want to share:

```bash
cd ~/Documents/my-project
~/Documents/devlink/devlink git serve .
```

Output:

```
Git daemon started for my-project.git (listening on 127.0.0.1:9418)
Git share ready! Teammates can connect via:
  devlink git connect abcd1234 my-project.git
```

* `abcd1234` is the temporary share token.
* `my-project.git` is the repo name teammates will use.
* Keep this terminal open while sharing.

---

### 2. Connect to a repository (Teammate machine)

Run:

```bash
./devlink git connect abcd1234 my-project.git
```

This opens a secure tunnel to the host’s Git repo.

---

### 3. Clone the repository

In a **new terminal** on the teammate’s machine:

```bash
git clone git://127.0.0.1:9418/my-project.git
```

This gives you a working copy of the repo.

---

### 4. Normal workflow

Once connected, teammates can use standard Git commands:

```bash
git add <file>
git commit -m "Your message"
git push origin master
git pull
```

---

## 📝 Notes

* You can serve **any repo** by `cd` into it and running:

  ```bash
  ./devlink git serve .
  ```
* If the repo has no commits yet, initialize it:

  ```bash
  git init
  git add .
  git commit -m "initial commit"
  ```
* Keep both `git serve` (host) and `git connect` (teammate) terminals open while working.
* When done, `Ctrl+C` to stop the tunnel.



```


