# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Model Context Protocol (MCP) server written in Go that provides a `search_photos` tool to query the Unsplash API. The server can run in two modes:
- **Stdio mode** (default): For direct MCP integrations via stdin/stdout
- **HTTP mode**: Exposes MCP streamable transport over HTTP

## Common Commands

### Development
```bash
# Run in stdio mode (default)
go run .

# Run HTTP server mode
go run . server -host 0.0.0.0 -port 8080

# Build the project
go build .
```

### Dependencies
Go automatically handles dependency management. Dependencies are declared in `go.mod` and will be fetched on first build/run.

## Architecture

### Core Components

- **`main.go`**: Entry point with command-line parsing, server setup, and transport handling
- **`unsplash.go`**: Unsplash API integration, data models, and the `search_photos` tool implementation
- **`logging_middleware.go`**: HTTP request/response logging for server mode

### Key Design Patterns

1. **Transport Abstraction**: Server supports both stdio and HTTP transports through the MCP SDK
2. **Tool Registration**: Tools are registered via `mcp.AddTool()` with handler functions
3. **Validation Layer**: Input validation with enum normalization and bounds checking
4. **Error Handling**: Structured error responses with context and timeout handling

### Data Flow

1. MCP client calls `search_photos` tool with parameters
2. `handleSearchPhotos()` validates and normalizes arguments
3. HTTP request made to Unsplash API with proper authentication
4. Response parsed and returned as structured JSON to MCP client

## Environment Requirements

- **Go 1.21+** (currently using 1.24.5)
- **UNSPLASH_ACCESS_KEY** environment variable (required)

## Tool: search_photos

The single exposed tool accepts these parameters:
- `query` (string, required): Search keyword
- `page` (number, optional): Page number (default: 1)
- `per_page` (number, optional): Results per page (default: 10, max: 30)
- `order_by` (string, optional): Sort order ("relevant" or "latest", default: "relevant")
- `color` (string, optional): Color filter (see allowedColors map in unsplash.go)
- `orientation` (string, optional): Orientation filter ("landscape", "portrait", "squarish")

## Dependencies

- `github.com/modelcontextprotocol/go-sdk`: Core MCP server functionality
- Standard library: HTTP client, JSON parsing, URL handling, context management

## Development Notes

- No tests are currently implemented
- Logging is built-in for HTTP mode via middleware
- Request timeout is hardcoded to 15 seconds
- Error messages include context for debugging
- All enum values are validated against predefined maps