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
	log.Printf("Creating MCP server...")

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in runHTTP: %v", r)
		}
	}()

	server := newUnsplashServer()
	log.Printf("MCP server created successfully")

	log.Printf("Creating streamable HTTP handler...")
	// 使用原始的NewStreamableHTTPHandler - 这是正确的API
	mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		log.Printf("Handler called for request: %s %s", req.Method, req.URL.Path)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in handler: %v", r)
			}
		}()
		return server
	}, nil)
	log.Printf("HTTP handler created successfully")

	// 创建一个多路复用器来处理不同的路径
	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"unsplash-mcp-server"}`))
	})

	// MCP端点
	mux.Handle("/", mcpHandler)

	log.Printf("Unsplash MCP server listening on http://%s", addr)
	log.Printf("Health check available at http://%s/health", addr)

	if err := http.ListenAndServe(addr, loggingHandler(mux)); err != nil {
		log.Printf("http server error: %v", err)
	}
}
