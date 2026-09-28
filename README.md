```text
  ____                  _ _    _____           _       __  __  ____ ____  
 | __ ) _   _ _ __   __| | |__|_   _|__   ___ | |     |  \/  |/ ___|  _ \ 
 |  _ \| | | | '_ \ / _` | / _ \| |/ _ \ / _ \| |_____| |\/| | |   | |_) |
 | |_) | |_| | | | | (_| | ||  /| | (_) | (_) | |_____| |  | | |___|  __/ 
 |____/ \__,_|_| |_|\__,_|_| \_\|_|\___/ \___/|_|     |_|  |_|\____|_|    
```

# BundleTool MCP Server

[![CI](https://github.com/iammohdzaki/bundletool-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/iammohdzaki/bundletool-mcp/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/iammohdzaki/bundletool-mcp?color=blue)](https://github.com/iammohdzaki/bundletool-mcp/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/iammohdzaki/bundletool-mcp)](https://golang.org/)
[![License](https://img.shields.io/github/license/iammohdzaki/bundletool-mcp)](LICENSE)

A fully-featured [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that empowers AI assistants (like Claude, Antigravity, and Cursor) to natively interact with Android's `bundletool`. 

This server allows LLMs to automatically build, extract, size, and install Android APKs directly to your connected devices from Android App Bundles (`.aab`).

## Features

- 📦 **Automated BundleTool Management**: Automatically downloads and manages the official Google `bundletool.jar` release in the background. No manual setup required!
- 🏗️ **`build_apks`**: Generate `.apks` archives from `.aab` bundles. Supports universal mode (with automatic `.apk` zip extraction), keystore signing, and connected-device targeting.
- 📱 **`install_apks`**: Directly deploy built APKs to connected devices or emulators via ADB. Supports `--allow-downgrade` and `--grant-runtime-permissions`.
- ✂️ **`extract_apks`**: Extract device-specific APKs using a hardware JSON spec.
- 📏 **`get_size`**: Calculate precise download sizes across different dimensions (ABI, SDK, Density).
- 🔍 **`get_device_spec`**: Inspect a connected device and generate a configuration JSON file.

## Prerequisites

The MCP server handles the `bundletool` JAR automatically, but your underlying OS must have:
1. **Java (JRE/JDK)**: Required to execute the JAR. Must be in your system `PATH`.
2. **Android SDK Platform-Tools (ADB)**: Required to communicate with connected devices. Must be in your system `PATH`.

## Installation & Configuration

Because this server interacts with your local Android hardware (ADB), it cannot be run inside an isolated Docker container. Choose one of the three installation methods below:

### Option 1: 1-Click Auto Installer (Easiest)
If you don't have Go installed, you can use our auto-installer. It will automatically download the correct binary for your OS from GitHub Releases and print the exact JSON configuration for you to paste into your AI editor.

**For Mac/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/iammohdzaki/bundletool-mcp/main/install.sh | bash
```

**For Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/iammohdzaki/bundletool-mcp/main/install.ps1 | iex
```

### Option 2: Zero-Install (`go run`)
If you are a developer with Go installed on your machine, you don't even need to download the binary manually. Tell your AI IDE to fetch and run the latest version directly from GitHub on the fly. 

Just copy and paste this into your MCP configuration (`~/.gemini/config/mcp_config.json` for Antigravity, or `claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "bundletool": {
      "command": "go",
      "args": ["run", "github.com/iammohdzaki/bundletool-mcp/cmd/bundletool-mcp@latest"]
    }
  }
}
```

### Option 3: Manual Binary Download
If you prefer complete control:
1. Go to the [GitHub Releases](https://github.com/iammohdzaki/bundletool-mcp/releases) page.
2. Download the binary for your specific OS and architecture.
3. Place it anywhere on your machine (e.g., `C:\mcp\bundletool-mcp.exe`).
4. Point your MCP config directly to the executable:
```json
{
  "mcpServers": {
    "bundletool": {
      "command": "C:\\mcp\\bundletool-mcp.exe",
      "args": []
    }
  }
}
```

## Security & Architecture

This server is built using the official [Anthropic Go SDK](https://github.com/modelcontextprotocol/go-sdk). 
It leverages `os/exec` to securely pass LLM-generated arguments directly to the OS binary layer without invoking an intermediate shell environment (like `bash` or `cmd.exe`). This structurally mitigates any risk of OS Command Injection.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
