# DevLink Configuration Guide

DevLink supports extensive configuration through environment variables for production-ready deployments.

## Security Configuration

All security issues identified in the audit have been resolved with configurable options.

### 🔒 Rate Limiting Configuration

Rate limiting is now fully configurable to handle production load requirements.

```bash
# Token generation rate limiting
export DEVLINK_TOKEN_RATE_LIMIT="100ms"    # Rate limit between token generations

# Validation rate limiting  
export DEVLINK_VALIDATION_MAX_ATTEMPTS="10"         # Max attempts per minute per client
export DEVLINK_VALIDATION_BLOCK_DURATION="5m"      # Block duration after exceeding limit
```

### 🧹 Cleanup Configuration

Cleanup intervals are configurable for optimal resource management.

```bash
# Token cleanup - how long codes remain valid
export DEVLINK_TOKEN_CLEANUP_INTERVAL="24h"

# Validation cleanup - how long to keep rate limit records
export DEVLINK_VALIDATION_CLEANUP_INTERVAL="1h"
```

### 🛡️ Offensive Words Configuration

The offensive words list is now loaded from external configuration files for better security.

```bash
# Path to offensive words configuration file
export DEVLINK_OFFENSIVE_WORDS_FILE="/path/to/offensive_words.txt"
```

**Example offensive words file format:**
```text
# Comments start with #
# One word per line
# Case insensitive
# Maximum 10,000 words
# Maximum 50 characters per word

test
demo
temp
fake
badword1
badword2
```

## Core Configuration

### Network Settings
```bash
export DEVLINK_P2P_PORT="8080"
export DEVLINK_P2P_NETWORK="devlink"
```

### Storage Settings
```bash
export DEVLINK_DATA_DIR="/path/to/data"
export DEVLINK_MAX_FILE_SIZE="1048576"  # 1MB in bytes
```

### Security Settings
```bash
export DEVLINK_ENCRYPTION_KEY="your-encryption-key"
export DEVLINK_DEFAULT_EXPIRY="1h"
```

### Logging Settings
```bash
export DEVLINK_LOG_LEVEL="info"         # debug, info, warn, error
export DEVLINK_LOG_FILE="/path/to/log"  # Optional log file
```

## Production Deployment Example

```bash
#!/bin/bash
# Production environment configuration

# Core settings
export DEVLINK_P2P_PORT="8080"
export DEVLINK_P2P_NETWORK="production"
export DEVLINK_DATA_DIR="/var/lib/devlink"
export DEVLINK_MAX_FILE_SIZE="2097152"  # 2MB
export DEVLINK_DEFAULT_EXPIRY="12h"

# Security settings
export DEVLINK_ENCRYPTION_KEY="${DEVLINK_ENCRYPTION_KEY}"  # From secure vault
export DEVLINK_OFFENSIVE_WORDS_FILE="/etc/devlink/offensive_words.txt"

# Rate limiting - production values
export DEVLINK_TOKEN_RATE_LIMIT="1s"                    # More conservative
export DEVLINK_VALIDATION_MAX_ATTEMPTS="5"             # Stricter limits
export DEVLINK_VALIDATION_BLOCK_DURATION="10m"         # Longer blocks

# Cleanup intervals - production optimized
export DEVLINK_TOKEN_CLEANUP_INTERVAL="12h"            # Shorter retention
export DEVLINK_VALIDATION_CLEANUP_INTERVAL="30m"       # More frequent cleanup

# Logging
export DEVLINK_LOG_LEVEL="warn"
export DEVLINK_LOG_FILE="/var/log/devlink/devlink.log"

# Start DevLink
./devlink
```

## Development Environment Example

```bash
#!/bin/bash
# Development environment configuration

# Relaxed rate limiting for testing
export DEVLINK_TOKEN_RATE_LIMIT="10ms"
export DEVLINK_VALIDATION_MAX_ATTEMPTS="50"
export DEVLINK_VALIDATION_BLOCK_DURATION="1m"

# Verbose logging
export DEVLINK_LOG_LEVEL="debug"

# Local storage
export DEVLINK_DATA_DIR="./data"

# Start DevLink
./devlink
```

## Configuration Validation

DevLink validates all configuration values on startup:

- **Duration formats**: Must be valid Go duration strings (e.g., "1h", "30m", "5s")
- **Numeric values**: Must be positive integers within valid ranges
- **File paths**: Must be accessible and secure (no path traversal)
- **Log levels**: Must be one of: debug, info, warn, error

Invalid configurations will cause DevLink to exit with a descriptive error message.

## Security Best Practices

1. **Use environment-specific configurations**
2. **Store encryption keys in secure vaults**
3. **Use stricter rate limits in production**
4. **Regularly update offensive words lists**
5. **Monitor logs for security events**
6. **Use shorter cleanup intervals in high-traffic environments**

## Defaults Reference

If environment variables are not set, DevLink uses these defaults:

```go
P2PPort:                      8080
P2PNetwork:                   "devlink"
DefaultExpiry:                "1h"
MaxFileSize:                  1048576  // 1MB
LogLevel:                     "info"
TokenRateLimit:               "100ms"
TokenCleanupInterval:         "24h"
ValidationMaxAttempts:        10
ValidationBlockDuration:      "5m"
ValidationCleanupInterval:    "1h"
```
