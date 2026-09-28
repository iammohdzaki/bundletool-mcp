# BundleTool MCP Server

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

## Installation

You can install the server binary directly using Go:

```bash
go install github.com/iammohdzaki/bundletool-mcp/cmd/bundletool-mcp@latest
```

*(This will place the `bundletool-mcp` binary in your `$GOPATH/bin` directory, which should be added to your system `PATH`.)*

## Usage (MCP Configuration)

To use this server with an MCP-compatible AI IDE or client, add the following to your MCP configuration file.

### For Antigravity (`~/.gemini/config/mcp_config.json`)
```json
{
  "mcpServers": {
    "bundletool": {
      "command": "bundletool-mcp",
      "args": []
    }
  }
}
```

### For Claude Desktop (`claude_desktop_config.json`)
```json
{
  "mcpServers": {
    "bundletool": {
      "command": "bundletool-mcp",
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
