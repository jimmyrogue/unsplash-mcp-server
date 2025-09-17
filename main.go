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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [stdio|server] [flags]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nModes:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  stdio   Run the MCP server over stdin/stdout (default)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  server  Run an HTTP server exposing the MCP tool\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nFlags (server mode only):\n")
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

func runStdio() {
	server := newUnsplashServer()
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("stdio server error: %v", err)
	}
}

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
