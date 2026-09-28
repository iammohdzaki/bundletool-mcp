#!/bin/bash
set -e

cat << "EOF"
  ____                  _ _    _____           _       __  __  ____ ____  
 | __ ) _   _ _ __   __| | |__|_   _|__   ___ | |     |  \/  |/ ___|  _ \ 
 |  _ \| | | | '_ \ / _` | / _ \| |/ _ \ / _ \| |_____| |\/| | |   | |_) |
 | |_) | |_| | | | | (_| | ||  /| | (_) | (_) | |_____| |  | | |___|  __/ 
 |____/ \__,_|_| |_|\__,_|_| \_\|_|\___/ \___/|_|     |_|  |_|\____|_|    
                                                                          
EOF
echo ""

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

BINARY_NAME="bundletool-mcp-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

INSTALL_DIR="$HOME/.bundletool-mcp/bin"
EXE_PATH="$INSTALL_DIR/bundletool-mcp"
if [ "$OS" = "windows" ]; then
    EXE_PATH="${EXE_PATH}.exe"
fi

echo "Fetching latest release for ${OS}-${ARCH}..."
DOWNLOAD_URL=$(curl -s https://api.github.com/repos/iammohdzaki/bundletool-mcp/releases/latest | grep "browser_download_url.*${BINARY_NAME}\"" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find a release binary for ${OS}-${ARCH}. Check GitHub releases."
    exit 1
fi

mkdir -p "$INSTALL_DIR"

echo "Downloading..."
curl -L "$DOWNLOAD_URL" -o "$EXE_PATH"
chmod +x "$EXE_PATH"

echo ""
echo "Successfully installed to $EXE_PATH"
echo ""
echo "======================================================="
echo "Please add the following to your MCP configuration JSON:"
echo "======================================================="
cat << EOF
{
  "mcpServers": {
    "bundletool": {
      "command": "$EXE_PATH",
      "args": []
    }
  }
}
EOF
