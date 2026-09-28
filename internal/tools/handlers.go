package tools

import (
	"bundletool-mcp/internal/bundletool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Handlers struct {
	client *bundletool.Client
}

func NewHandlers(client *bundletool.Client) *Handlers {
	return &Handlers{client: client}
}

// Register maps our Go structs to JSON Schemas and binds them to the server.
func (h *Handlers) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "build_apks",
		Description: "Converts an Android App Bundle (.aab) into an archive of APKs (.apks).",
	}, h.handleBuildApks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "install_apks",
		Description: "Installs an .apks file onto a connected device/emulator via ADB.",
	}, h.handleInstallApks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_size",
		Description: "Calculates the total size of the APKs targeting a specific device configuration.",
	}, h.handleGetSize)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_device_spec",
		Description: "Generates a JSON file containing the hardware/software specifications of a connected device.",
	}, h.handleGetDeviceSpec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "extract_apks",
		Description: "Extracts from an APK Set the APKs that should be installed on a given device based on a device spec JSON.",
	}, h.handleExtractApks)
}

// formatResult maps a successful shell output or error into an MCP response.
func formatResult(out string, err error) (*mcp.CallToolResult, any, error) {
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
		}, nil, nil
	}
	if out == "" {
		out = "Success."
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: out},
		},
	}, nil, nil
}
