package tools

import (
	"context"
	"fmt"

	"bundletool-mcp/internal/bundletool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BuildApksArgs defines the JSON schema arguments required by the LLM to call the build_apks tool.
type BuildApksArgs struct {
	AabPath         string `json:"aab_path" jsonschema:"Absolute path to the input .aab file."`
	ApksPath        string `json:"apks_path" jsonschema:"Absolute path where the .apks file should be saved."`
	Overwrite       bool   `json:"overwrite,omitempty" jsonschema:"If true, overwrites any existing output file."`
	Aapt2           string `json:"aapt2,omitempty" jsonschema:"Custom path to AAPT2."`
	KsPath          string `json:"ks_path,omitempty" jsonschema:"Path to the deployment keystore for signing."`
	KsPass          string `json:"ks_pass,omitempty" jsonschema:"Keystore password (e.g. pass:qwerty). Required if ks_path is used."`
	KsKeyAlias      string `json:"ks_key_alias,omitempty" jsonschema:"Alias of the signing key. Required if ks_path is used."`
	KeyPass         string `json:"key_pass,omitempty" jsonschema:"Password for the signing key. Required if ks_path is used."`
	ConnectedDevice bool   `json:"connected_device,omitempty" jsonschema:"Build APKs targeting a connected device configuration."`
	DeviceID        string `json:"device_id,omitempty" jsonschema:"Serial ID of the connected device to target."`
	DeviceSpec      string `json:"device_spec,omitempty" jsonschema:"Path to a device spec JSON file."`
	Mode            string `json:"mode,omitempty" jsonschema:"Mode: default, universal, system, persistent, instant, archive."`
	LocalTesting    bool   `json:"local_testing,omitempty" jsonschema:"If true, builds in local testing mode."`
	Modules         string `json:"modules,omitempty" jsonschema:"Comma-separated list of modules to include."`
}

// handleBuildApks is the MCP tool handler for 'build_apks'.
// It validates keystore dependencies and delegates to the bundletool client.
func (h *Handlers) handleBuildApks(ctx context.Context, req *mcp.CallToolRequest, args BuildApksArgs) (*mcp.CallToolResult, any, error) {
	if args.KsPath != "" {
		if args.KsPass == "" || args.KsKeyAlias == "" || args.KeyPass == "" {
			return formatResult("", fmt.Errorf("validation error: ks_pass, ks_key_alias, and key_pass are required when ks_path is provided (non-interactive MCP cannot prompt for passwords)"))
		}
	}

	opts := bundletool.BuildApksOptions{
		AabPath:         args.AabPath,
		ApksPath:        args.ApksPath,
		Overwrite:       args.Overwrite,
		Aapt2:           args.Aapt2,
		KsPath:          args.KsPath,
		KsPass:          args.KsPass,
		KsKeyAlias:      args.KsKeyAlias,
		KeyPass:         args.KeyPass,
		ConnectedDevice: args.ConnectedDevice,
		DeviceID:        args.DeviceID,
		DeviceSpec:      args.DeviceSpec,
		Mode:            args.Mode,
		LocalTesting:    args.LocalTesting,
		Modules:         args.Modules,
	}
	out, err := h.client.BuildApks(ctx, opts)
	return formatResult(out, err)
}
