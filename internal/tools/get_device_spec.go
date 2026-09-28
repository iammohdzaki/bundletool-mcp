package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetDeviceSpecArgs struct {
	OutputPath string `json:"output_path" jsonschema:"Absolute path where the JSON file will be saved."`
	DeviceID   string `json:"device_id,omitempty" jsonschema:"Specific ADB device serial number."`
}

func (h *Handlers) handleGetDeviceSpec(ctx context.Context, req *mcp.CallToolRequest, args GetDeviceSpecArgs) (*mcp.CallToolResult, any, error) {
	out, err := h.client.GetDeviceSpec(ctx, args.OutputPath, args.DeviceID)
	return formatResult(out, err)
}
