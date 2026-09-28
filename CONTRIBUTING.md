# Contributing to BundleTool MCP

First off, thank you for considering contributing to this project! 

## Getting Started

1. Fork the repository and clone it locally.
2. Ensure you have Go 1.22+ installed.
3. Run `go mod download` to fetch dependencies.
4. Verify your environment has Java and ADB installed.

## Project Structure

The project is structured using standard Go architecture:
- `cmd/bundletool-mcp/`: Contains the main entrypoint and dependency injection wiring.
- `internal/bundletool/`: The core business logic for wrapping the `bundletool.jar` CLI.
- `internal/executor/`: Shell execution abstractions (allows for easy unit testing).
- `internal/tools/`: The MCP tool definitions and JSON schema mappings.

## Adding a New Tool

1. Define the input arguments struct in `internal/tools/your_tool.go` using the `jsonschema` struct tags.
2. Add the business logic to `internal/bundletool/client.go`.
3. Create the handler function in your new tool file.
4. Register the tool in `internal/tools/handlers.go` under the `Register` function.

## Pull Requests

1. Create a feature branch.
2. Make sure your code is formatted with `go fmt ./...`.
3. Ensure all tests pass with `go test ./...`.
4. Submit a PR describing your changes in detail!
