package bundletool

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// MockExecutor implements the executor.Executor interface for testing
type MockExecutor struct {
	LastCommand string
	LastArgs    []string
	ShouldFail  bool
}

func (m *MockExecutor) Execute(ctx context.Context, command string, args ...string) (string, error) {
	m.LastCommand = command
	m.LastArgs = args

	if m.ShouldFail {
		return "", fmt.Errorf("mock error")
	}

	// Simulate bundletool success
	return "Mock execution success", nil
}

func (m *MockExecutor) LookPath(file string) (string, error) {
	return "/usr/bin/" + file, nil
}

func TestBuildApks(t *testing.T) {
	// 1. Setup our mock dependencies
	mockExec := &MockExecutor{}
	client := NewClient(mockExec, "/fake/path/bundletool.jar")

	// 2. Call the method we want to test
	ctx := context.Background()
	opts := BuildApksOptions{
		AabPath:    "app.aab",
		ApksPath:   "app.apks",
		KsPath:     "my.keystore",
		KsPass:     "pass:1234",
		KsKeyAlias: "key0",
		KeyPass:    "pass:1234",
		Mode:       "universal",
	}
	_, err := client.BuildApks(ctx, opts)

	// 3. Assert the results
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ensure it tried to run 'java'
	if mockExec.LastCommand != "java" {
		t.Errorf("Expected command to be 'java', got %s", mockExec.LastCommand)
	}

	// Ensure the arguments were constructed correctly
	expectedArgs := "-jar /fake/path/bundletool.jar build-apks --bundle=app.aab --output=app.apks --ks=my.keystore --ks-pass=pass:1234 --ks-key-alias=key0 --key-pass=pass:1234 --mode=universal"
	actualArgs := strings.Join(mockExec.LastArgs, " ")

	if actualArgs != expectedArgs {
		t.Errorf("Expected args:\n%s\nGot:\n%s", expectedArgs, actualArgs)
	}
}
