package tools

import (
	"context"

	"bundletool-mcp/internal/bundletool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetSizeArgs struct {
	ApksPath   string `json:"apks_path" jsonschema:"Absolute path to the .apks file."`
	DeviceSpec string `json:"device_spec,omitempty" jsonschema:"Path to the device spec JSON file."`
	Dimensions string `json:"dimensions,omitempty" jsonschema:"Dimensions for size estimates (e.g. SDK, ABI, SCREEN_DENSITY, LANGUAGE, ALL)."`
	Modules    string `json:"modules,omitempty" jsonschema:"Comma-separated list of modules to consider."`
	Instant    bool   `json:"instant,omitempty" jsonschema:"Measure download size of instant-enabled APKs instead."`
}

func (h *Handlers) handleGetSize(ctx context.Context, req *mcp.CallToolRequest, args GetSizeArgs) (*mcp.CallToolResult, any, error) {
	opts := bundletool.GetSizeOptions{
		ApksPath:   args.ApksPath,
		DeviceSpec: args.DeviceSpec,
		Dimensions: args.Dimensions,
		Modules:    args.Modules,
		Instant:    args.Instant,
	}
	out, err := h.client.GetSize(ctx, opts)
	return formatResult(out, err)
}
