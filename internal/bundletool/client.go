package bundletool

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"bundletool-mcp/internal/executor"
)

// Client handles all interaction with the bundletool JAR.
// It is completely independent of the MCP protocol.
type Client struct {
	executor executor.Executor
	jarPath  string
}

func NewClient(exec executor.Executor, jarPath string) *Client {
	return &Client{
		executor: exec,
		jarPath:  jarPath,
	}
}

func (c *Client) ensureJava() error {
	if _, err := c.executor.LookPath("java"); err != nil {
		return fmt.Errorf("Java is not installed or not in the system PATH. BundleTool requires Java to run. Please inform the user to install a Java Runtime Environment (JRE/JDK).")
	}
	return nil
}

func (c *Client) ensureAdbAndDevice(ctx context.Context, deviceID string) error {
	if _, err := c.executor.LookPath("adb"); err != nil {
		return fmt.Errorf("ADB is not installed or not in the system PATH. Please inform the user they need to install Android SDK Platform-Tools and add it to their environment variables")
	}

	// Check if a device is actually connected
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "get-state")

	out, err := c.executor.Execute(ctx, "adb", args...)
	if err != nil {
		return fmt.Errorf("No Android device is currently connected (or it is unauthorized). Please ask the user to plug in their device or start an emulator. ADB Output: %s", out)
	}

	return nil
}

type BuildApksOptions struct {
	AabPath         string
	ApksPath        string
	Overwrite       bool
	Aapt2           string
	KsPath          string
	KsPass          string
	KsKeyAlias      string
	KeyPass         string
	ConnectedDevice bool
	DeviceID        string
	DeviceSpec      string
	Mode            string
	LocalTesting    bool
	Modules         string
}

func (c *Client) BuildApks(ctx context.Context, opts BuildApksOptions) (string, error) {
	if err := c.ensureJava(); err != nil {
		return "", err
	}
	if opts.ConnectedDevice {
		if err := c.ensureAdbAndDevice(ctx, opts.DeviceID); err != nil {
			return "", err
		}
	}

	args := []string{"-jar", c.jarPath, "build-apks", "--bundle=" + opts.AabPath, "--output=" + opts.ApksPath}

	if opts.Overwrite {
		args = append(args, "--overwrite")
	}
	if opts.Aapt2 != "" {
		args = append(args, "--aapt2="+opts.Aapt2)
	}
	if opts.KsPath != "" {
		args = append(args, "--ks="+opts.KsPath)
	}
	if opts.KsPass != "" {
		args = append(args, "--ks-pass="+opts.KsPass)
	}
	if opts.KsKeyAlias != "" {
		args = append(args, "--ks-key-alias="+opts.KsKeyAlias)
	}
	if opts.KeyPass != "" {
		args = append(args, "--key-pass="+opts.KeyPass)
	}
	if opts.ConnectedDevice {
		args = append(args, "--connected-device")
	}
	if opts.DeviceID != "" {
		args = append(args, "--device-id="+opts.DeviceID)
	}
	if opts.DeviceSpec != "" {
		args = append(args, "--device-spec="+opts.DeviceSpec)
	}
	if opts.Mode != "" {
		args = append(args, "--mode="+opts.Mode)
	}
	if opts.LocalTesting {
		args = append(args, "--local-testing")
	}
	if opts.Modules != "" {
		args = append(args, "--modules="+opts.Modules)
	}

	out, err := c.executor.Execute(ctx, "java", args...)
	if err != nil {
		return out, err
	}

	if opts.Mode == "universal" {
		err := extractUniversalApk(opts.ApksPath)
		if err != nil {
			out += fmt.Sprintf("\nWarning: Built successfully, but failed to extract universal.apk: %v", err)
		} else {
			out += fmt.Sprintf("\nSuccessfully extracted universal.apk to the same directory as %s", opts.ApksPath)
		}
	}

	return out, nil
}

func extractUniversalApk(apksPath string) error {
	r, err := zip.OpenReader(apksPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "universal.apk" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			destPath := filepath.Join(filepath.Dir(apksPath), "universal.apk")
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err = io.Copy(outFile, rc); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("universal.apk not found inside archive")
}

type InstallApksOptions struct {
	ApksPath                string
	DeviceID                string
	Modules                 string
	AllowDowngrade          bool
	AllowTestOnly           bool
	GrantRuntimePermissions bool
}

func (c *Client) InstallApks(ctx context.Context, opts InstallApksOptions) (string, error) {
	if err := c.ensureJava(); err != nil {
		return "", err
	}
	if err := c.ensureAdbAndDevice(ctx, opts.DeviceID); err != nil {
		return "", err
	}

	args := []string{"-jar", c.jarPath, "install-apks", "--apks=" + opts.ApksPath}

	if opts.DeviceID != "" {
		args = append(args, "--device-id="+opts.DeviceID)
	}
	if opts.Modules != "" {
		args = append(args, "--modules="+opts.Modules)
	}
	if opts.AllowDowngrade {
		args = append(args, "--allow-downgrade")
	}
	if opts.AllowTestOnly {
		args = append(args, "--allow-test-only")
	}
	if opts.GrantRuntimePermissions {
		args = append(args, "--grant-runtime-permissions")
	}

	return c.executor.Execute(ctx, "java", args...)
}

type GetSizeOptions struct {
	ApksPath   string
	DeviceSpec string
	Dimensions string
	Modules    string
	Instant    bool
}

func (c *Client) GetSize(ctx context.Context, opts GetSizeOptions) (string, error) {
	if err := c.ensureJava(); err != nil {
		return "", err
	}

	args := []string{"-jar", c.jarPath, "get-size", "total", "--apks=" + opts.ApksPath}

	if opts.DeviceSpec != "" {
		args = append(args, "--device-spec="+opts.DeviceSpec)
	}
	if opts.Dimensions != "" {
		args = append(args, "--dimensions="+opts.Dimensions)
	}
	if opts.Modules != "" {
		args = append(args, "--modules="+opts.Modules)
	}
	if opts.Instant {
		args = append(args, "--instant")
	}

	return c.executor.Execute(ctx, "java", args...)
}

type ExtractApksOptions struct {
	ApksPath        string
	DeviceSpec      string
	OutputDir       string
	Modules         string
	Instant         bool
	IncludeMetadata bool
}

func (c *Client) ExtractApks(ctx context.Context, opts ExtractApksOptions) (string, error) {
	if err := c.ensureJava(); err != nil {
		return "", err
	}

	args := []string{"-jar", c.jarPath, "extract-apks", "--apks=" + opts.ApksPath, "--device-spec=" + opts.DeviceSpec}

	if opts.OutputDir != "" {
		args = append(args, "--output-dir="+opts.OutputDir)
	}
	if opts.Modules != "" {
		args = append(args, "--modules="+opts.Modules)
	}
	if opts.Instant {
		args = append(args, "--instant")
	}
	if opts.IncludeMetadata {
		args = append(args, "--include-metadata")
	}

	return c.executor.Execute(ctx, "java", args...)
}

func (c *Client) GetDeviceSpec(ctx context.Context, outputPath, deviceID string) (string, error) {
	if err := c.ensureJava(); err != nil {
		return "", err
	}
	if err := c.ensureAdbAndDevice(ctx, deviceID); err != nil {
		return "", err
	}

	args := []string{"-jar", c.jarPath, "get-device-spec", "--output=" + outputPath}

	if deviceID != "" {
		args = append(args, "--device-id="+deviceID)
	}

	return c.executor.Execute(ctx, "java", args...)
}
