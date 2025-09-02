# devlink

## Folder Structure

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


## **Project Structure**

```
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
```

---


## **Installation**

1. Clone the DevLink repository:

```bash
git clone <devlink-repo-url>
cd DevLink
```

2. Build the CLI tool:

```bash
go build -o devlink
```

* This generates the `devlink` executable.

---

## **Git Repository Sharing**

DevLink allows you to **share a local Git repository** with teammates over a secure tunnel without exposing your machine publicly.

### **1. Serve a Repository (Laptop A / Server)**

```bash
# Initialize a bare repository
git init --bare ~/test-repo.git
touch ~/test-repo.git/git-daemon-export-ok

# Start DevLink Git server
./devlink git serve ~/test-repo.git
```

* Output will show the DevLink token to share:

```
Git share ready! Teammates can connect via:
  devlink git connect <token> 9418
```

* Keep this terminal running.

### **2. Connect to a Repository (Laptop B / Client)**

```bash
# Connect to the shared repo via the DevLink token
./devlink git connect <token> test-repo.git
```

* Keep this terminal open to maintain the tunnel.

### **3. Clone the Repository**

In another terminal on Laptop B:

```bash
git clone git://127.0.0.1:9418/test-repo.git
```

* Use a new folder if the destination folder already exists:

```bash
git clone git://127.0.0.1:9418/test-repo.git my-clone
```

### **4. First Commit (if repo is empty)**

On Laptop A (or after cloning locally):

```bash
# Clone bare repo to a temporary working directory
git clone ~/test-repo.git ~/temp-clone
cd ~/temp-clone
touch README.md
git add README.md
git commit -m "Initial commit"
git push origin master
```

---

## **Normal Workflow**

* After the initial setup, you can use normal Git commands through the DevLink tunnel:

```bash
git add <file>
git commit -m "message"
git push origin master
git pull
```

* Keep `git serve` running on the server and `git connect` running on the client.

