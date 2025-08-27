# DevLink

A CLI tool for development workflow management with smart features for environment file sharing.

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

### 📋 Smart Features

#### 🎯 Auto-completion
Get intelligent suggestions for share codes, file paths, and commands:

```bash
# Generate shell completion scripts
devlink env completion bash > ~/.bash_completion
devlink env completion zsh > ~/.zsh_completion

# Get suggestions for share codes
devlink env completion suggest sharecode "blue-whale"

# Get file suggestions
devlink env completion suggest file ".env"
```

#### 📋 Templates
Use pre-built environment templates for common development scenarios:

```bash
# List available templates
devlink env template list

# Show template details
devlink env template show nodejs

# Create environment file from template
devlink env template create nodejs --output .env

# Search templates
devlink env template search backend
```

#### 📦 Bulk Operations
Share multiple environment files at once:

```bash
# Share multiple files
devlink env bulk share file1.env file2.env file3.env

# Share with prefix and custom expiry
devlink env bulk share *.env --prefix myproject --expiry 24h

# Share with grouping
devlink env bulk share .env* --group-by category
```

#### 🔍 Search & Filter
Find specific variables across multiple environment files:

```bash
# Search for variables containing "DATABASE"
devlink env search "DATABASE"

# Show only sensitive variables
devlink env search --sensitive

# Search by category
devlink env search --categories database,api

# Use regex patterns
devlink env search --regex "API_.*"

# Get variable suggestions
devlink env search suggest "DAT"
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
│       ├── get.go       # devlink env get
│       ├── template.go  # Template management
│       ├── bulk.go      # Bulk operations
│       ├── search.go    # Search & filter
│       └── completion.go # Auto-completion
├── internal/
│   ├── env/             # Environment processing
│   │   ├── parser.go    # .env file parser
│   │   ├── validator.go # Security validation
│   │   ├── formatter.go # Output formatting
│   │   ├── bulk.go      # Bulk operations
│   │   └── server/      # HTTP server
│   └── util/
│       ├── config.go    # Configuration
│       ├── templates.go # Template management
│       ├── search.go    # Search functionality
│       ├── completion.go # Auto-completion
│       ├── token.go     # Share code generation
│       └── encryption.go # Data encryption
├── main.go              # Entry point
├── go.mod
└── go.sum
```
