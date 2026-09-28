package tools

import (
	"context"

	"bundletool-mcp/internal/bundletool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InstallApksArgs struct {
	ApksPath                string `json:"apks_path" jsonschema:"Absolute path to the generated .apks file."`
	DeviceID                string `json:"device_id,omitempty" jsonschema:"Specific ADB device serial number."`
	Modules                 string `json:"modules,omitempty" jsonschema:"Comma-separated list of modules to install."`
	AllowDowngrade          bool   `json:"allow_downgrade,omitempty" jsonschema:"Allow install even if app has lower version code."`
	AllowTestOnly           bool   `json:"allow_test_only,omitempty" jsonschema:"Allow deployment of apps with android:testOnly=true."`
	GrantRuntimePermissions bool   `json:"grant_runtime_permissions,omitempty" jsonschema:"Grant runtime permissions automatically (Android M+)."`
}

func (h *Handlers) handleInstallApks(ctx context.Context, req *mcp.CallToolRequest, args InstallApksArgs) (*mcp.CallToolResult, any, error) {
	opts := bundletool.InstallApksOptions{
		ApksPath:                args.ApksPath,
		DeviceID:                args.DeviceID,
		Modules:                 args.Modules,
		AllowDowngrade:          args.AllowDowngrade,
		AllowTestOnly:           args.AllowTestOnly,
		GrantRuntimePermissions: args.GrantRuntimePermissions,
	}
	out, err := h.client.InstallApks(ctx, opts)
	return formatResult(out, err)
}
