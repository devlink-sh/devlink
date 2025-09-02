
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

## **Installation**

1. **Clone the DevLink repository**

```bash
git clone <devlink-repo-url>
cd DevLink
```

2. **Build the CLI tool**

```bash
go build -o devlink
```

> This generates the `devlink` executable in your current directory.

---

## **Git Repository Sharing**

DevLink allows you to securely share a local Git repository with teammates via a temporary, encrypted tunnel.

---

### **1. Serve a Repository (Host Machine / Laptop A)**

1. Initialize a bare repository:

```bash
git init --bare ~/test-repo.git
touch ~/test-repo.git/git-daemon-export-ok
```

2. Start the DevLink Git server:

```bash
./devlink git serve ~/test-repo.git
```

> The terminal will display a token for connecting teammates:

```
Git share ready! Teammates can connect via:
  devlink git connect <token> 9418
```

> Keep this terminal open to maintain the server session.

---

### **2. Connect to a Repository (Client Machine / Laptop B)**

```bash
./devlink git connect <token> test-repo.git
```

> Keep this terminal open to maintain the tunnel.

---

### **3. Clone the Repository**

In a separate terminal on Laptop B:

```bash
git clone git://127.0.0.1:9418/test-repo.git
```

> If the target folder already exists, specify a new folder name:

```bash
git clone git://127.0.0.1:9418/test-repo.git my-clone
```

---

### **4. Initial Commit (if repository is empty)**

On Laptop A (or after cloning locally):

```bash
git clone ~/test-repo.git ~/temp-clone
cd ~/temp-clone
touch README.md
git add README.md
git commit -m "Initial commit"
git push origin master
```

> This ensures the repository has a starting commit for collaboration.

---

## **Normal Workflow**

Once the repository is set up, you can use standard Git commands through the DevLink tunnel:

```bash
git add <file>
git commit -m "Your message"
git push origin master
git pull
```

> **Important:**
>
> * Keep `git serve` running on the host machine.
> * Keep `git connect` running on the client machine to maintain the tunnel.

```


