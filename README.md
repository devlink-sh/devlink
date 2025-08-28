# DevLink 🔐

> **"Sharing secrets safely, one .env file at a time!"** 🚀

Ever tried to share your environment variables with a teammate and thought, *"Hmm, maybe I shouldn't paste this API key in Slack..."*? 

Well, worry no more! DevLink is here to save the day with **zero-trust networking** that makes sharing .env files as safe as passing a secret note through an invisible, encrypted tunnel. ✨

## 🌟 What Makes DevLink Special?

- 🔐 **Zero-trust networking** - No internet exposure, no worries!
- 🛡️ **End-to-end encryption** - Your secrets stay secret
- 🔍 **Smart security detection** - Automatically finds and protects sensitive data
- 📤 **One-click sharing** - Share with a simple code (like `blue-dragon-123`)
- 💥 **Self-destructing shares** - Files disappear after use (Mission Impossible style!)
- 🎯 **Beginner-friendly** - No PhD in cryptography required

## 🚀 Quick Start (5 minutes to awesome!)

### Step 1: Build the Magic
```bash
go build -o devlink
```

### Step 2: Start the Secret Tunnel
```bash
./devlink server
```
*Keep this running while you want to share files!*

### Step 3: Share Your .env File
```bash
./devlink env share .env --expiry 1h
```
*This gives you a code like `blue-dragon-123` to share with your teammate*

### Step 4: Your Teammate Gets the File
```bash
./devlink env get blue-dragon-123 --output .env
```
*Poof! The .env file appears like magic! ✨*

## 📚 Commands Made Simple

### 🚀 Server Commands
```bash
./devlink server                    # Start the secure tunnel
./devlink server --service my-team  # Custom service name
./devlink server --verbose          # See the magic happening
```

### 🔐 Environment Sharing
```bash
# Share your .env file
./devlink env share .env                    # Share for 1 hour
./devlink env share .env --expiry 24h       # Share for 24 hours
./devlink env share .env --readonly         # Make it read-only (safer!)

# Get a shared .env file
./devlink env get blue-dragon-123          # Get and display
./devlink env get blue-dragon-123 --output .env  # Save to file
./devlink env get blue-dragon-123 --unmask       # Show secrets (be careful!)
```

## 🎯 Perfect For...

- **Development teams** sharing environment configurations
- **DevOps engineers** distributing secrets safely
- **Anyone** who's tired of Slack DMs with API keys
- **Security-conscious developers** who want zero-trust networking
- **People** who like cool CLI tools with emojis! 🎉

## 🔧 How It Works (The Magic Explained)

1. **You share a file** → DevLink encrypts it and creates a temporary code
2. **Your teammate uses the code** → DevLink decrypts and delivers the file
3. **The share disappears** → Like a self-destructing message! 💥
4. **Zero internet exposure** → Everything goes through secure tunnels

Think of it like having a secret handshake that only your team knows, but for files! 🤝

## 🛡️ Security Features

- **Zero-trust networking** - No network exposure, ever
- **Identity-based access** - Only your team can connect
- **Encrypted tunnels** - All communication is encrypted
- **Sensitive data masking** - Automatically detects and protects secrets
- **Single-use shares** - Files delete themselves after access
- **Time-based expiration** - Shares expire automatically

## 🎨 Architecture (For the Curious)

```
pkg/envsharing/
├── core/           # The brain 🧠
│   ├── parser.go   # Reads .env files
│   ├── validator.go # Checks for security issues
│   ├── formatter.go # Makes output pretty
│   ├── types.go    # Data structures
│   └── encryption/ # The secret sauce 🔐
├── network/        # The tunnel 🌉
│   ├── ziti.go     # OpenZiti service
│   └── client.go   # OpenZiti client
└── cli/            # The friendly face 😊
    └── commands.go # User commands
```

## ⚙️ Configuration (Optional)

Set these environment variables if you want to customize:

```bash
export ZITI_CONTROLLER_URL="https://controller.example.com"
export ZITI_IDENTITY_FILE="/path/to/identity.json"
export ZITI_SERVICE_NAME="my-team-service"
```

## 📄 License

MIT License - Feel free to use this in your projects!

---

**Made with ❤️ for developers who care about security and convenience!**

*"Because sharing should be caring, not scary!"* 🎭

