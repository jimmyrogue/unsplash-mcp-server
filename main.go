package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	host = flag.String("host", "127.0.0.1", "host to listen on")
	port = flag.Int("port", 8080, "port to listen on")
)

const usageTemplate = `Usage: %s [stdio|server] [flags]

Modes:
  stdio   Run the MCP server over stdin/stdout (default)
  server  Run an HTTP server exposing the MCP tool

Flags (server mode only):
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageTemplate, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	mode := "stdio"
	if flag.NArg() > 0 {
		mode = flag.Arg(0)
	}

	switch mode {
	case "stdio":
		runStdio()
	case "server":
		runHTTP(fmt.Sprintf("%s:%d", *host, *port))
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n\n", mode)
		flag.Usage()
		os.Exit(2)
	}
}

// newUnsplashServer wires up the MCP server instance with all Unsplash tools.
func newUnsplashServer() *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "Unsplash MCP Server",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search_photos",
		Description: "Search for Unsplash photos",
	}, handleSearchPhotos)

	return srv
}

// runStdio serves the MCP server over stdin/stdout for local tool integrations.
func runStdio() {
	server := newUnsplashServer()
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("stdio server error: %v", err)
	}
}

// runHTTP exposes the server through the MCP streamable HTTP transport.
func runHTTP(addr string) {
	server := newUnsplashServer()
	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)

	log.Printf("Unsplash MCP server listening on http://%s", addr)

	if err := http.ListenAndServe(addr, loggingHandler(handler)); err != nil {
		log.Printf("http server error: %v", err)
	}
}
