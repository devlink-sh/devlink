# DevLink CLI Installation

## Quick Install (Linux)

Install DevLink CLI globally on your Linux system with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/knk/devlink/main/install.sh | bash
```

Or if you prefer wget:

```bash
wget -qO- https://raw.githubusercontent.com/knk/devlink/main/install.sh | bash
```

## Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/knk/devlink.git
cd devlink
```

2. Run the installer:
```bash
./install.sh
```

## What it does

- Downloads the appropriate binary for your Linux architecture (amd64, arm64, arm)
- Installs it globally to `/usr/local/bin/`
- Makes the `devlink` command available from anywhere in your system
- Automatically detects your system architecture
- Includes safety checks and error handling

## Requirements

- Linux system (x86_64, ARM64, or ARM)
- curl or wget
- sudo access (for global installation)

## After Installation

Verify the installation:
```bash
devlink --version
```

Get help:
```bash
devlink --help
```

## Uninstall

To remove DevLink CLI:
```bash
sudo rm /usr/local/bin/devlink
```
