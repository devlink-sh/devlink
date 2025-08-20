# DevLink

A CLI tool for development workflow management.

## Installation

```bash
go build -o devlink
```

## Usage

### 🔐 Environment File Sharing

Share environment files securely with your team:

```bash
# Share an environment file
devlink env share .env

# Share with custom expiry and read-only flag
devlink env share .env --expiry 24h --readonly

# Get a shared environment file
devlink env get ABC123

# Save retrieved file to disk
devlink env get ABC123 --output .env
```

### Example Output

**Sharing a file:**
```
🚀 Sharing: .env
⏰ Expires: 1h

✨ Share created successfully!
📋 Share this code with your team:
   ABC123

💡 Use: devlink env get ABC123
```

**Getting a file (example content):**
```
🔍 Retrieving: ABC123

📄 Environment file content:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=your-secret-key
REDIS_URL=redis://localhost:6379
NODE_ENV=development
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

*Note: The environment file content above is just an example. Your actual file content will be displayed.*

## Project Structure

```
devlink/
├── cmd/
│   ├── root.go          # Root CLI command
│   └── env/             # Environment commands
│       ├── env.go       # Main env command
│       ├── share.go     # devlink env share
│       └── get.go       # devlink env get
├── internal/
│   └── util/
│       └── config.go    # Configuration
├── main.go              # Entry point
├── go.mod
└── go.sum
```
