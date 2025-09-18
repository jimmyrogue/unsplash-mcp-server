#!/usr/bin/env node

import { config } from "dotenv";
import express from "express";
import { randomUUID } from "node:crypto";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";
import { createMcpServer } from "./server.js";
import { ServerConfig } from "./types/index.js";
import { createLoggingMiddleware } from "./middleware/logging.js";

config();

const SERVER_CONFIG: ServerConfig = {
  name: "Unsplash MCP Server",
  version: "1.0.0",
  host: process.env.HOST || "127.0.0.1",
  port: parseInt(process.env.PORT || "8080", 10),
};

function showUsage() {
  console.error(`Usage: ${process.argv[1]} [stdio|server] [flags]

Modes:
  stdio   Run the MCP server over stdin/stdout (default)
  server  Run an HTTP server exposing the MCP tool

Environment Variables:
  UNSPLASH_ACCESS_KEY  Your Unsplash API access key (required)
  HOST                 Host to bind to (default: 127.0.0.1)
  PORT                 Port to listen on (default: 8080)
`);
}

async function runStdio() {
  const server = createMcpServer(SERVER_CONFIG);

  const transport = new StdioServerTransport();
  await server.connect(transport);
}

async function runHttp(host: string, port: number) {
  console.log("Creating MCP server...");

  const app = express();
  app.use(express.json());
  app.use(createLoggingMiddleware({ logRequests: true, logResponses: true }));

  // Map to store transports by session ID
  const transports: Record<string, StreamableHTTPServerTransport> = {};

  // Health check endpoint
  app.get("/health", (_req, res) => {
    res.json({
      status: "ok",
      service: "unsplash-mcp-server-ts",
      version: SERVER_CONFIG.version,
      timestamp: new Date().toISOString()
    });
  });

  // Handle POST requests for client-to-server communication
  app.post("/mcp", async (req, res) => {
    const sessionId = req.headers["mcp-session-id"] as string | undefined;
    let transport: StreamableHTTPServerTransport;

    try {
      if (sessionId && transports[sessionId]) {
        transport = transports[sessionId];
      } else {
        transport = new StreamableHTTPServerTransport({
          sessionIdGenerator: () => randomUUID(),
          onsessioninitialized: (sessionId) => {
            transports[sessionId] = transport;
          },
        });

        transport.onclose = () => {
          if (transport.sessionId) {
            delete transports[transport.sessionId];
          }
        };

        const server = createMcpServer(SERVER_CONFIG);

        await server.connect(transport);
      }

      await transport.handleRequest(req, res, req.body);
    } catch (error) {
      console.error("Error handling MCP request:", error);
      if (!res.headersSent) {
        res.status(500).json({
          jsonrpc: "2.0",
          error: {
            code: -32603,
            message: "Internal server error",
          },
          id: null,
        });
      }
    }
  });

  // Handle GET requests for server-to-client notifications via SSE
  app.get("/mcp", async (req, res) => {
    const sessionId = req.headers["mcp-session-id"] as string | undefined;
    if (!sessionId || !transports[sessionId]) {
      res.status(400).send("Invalid or missing session ID");
      return;
    }

    const transport = transports[sessionId];
    await transport.handleRequest(req, res);
  });

  // Handle DELETE requests for session termination
  app.delete("/mcp", async (req, res) => {
    const sessionId = req.headers["mcp-session-id"] as string | undefined;
    if (!sessionId || !transports[sessionId]) {
      res.status(400).send("Invalid or missing session ID");
      return;
    }

    const transport = transports[sessionId];
    await transport.handleRequest(req, res);
  });

  app.listen(port, host, () => {
    console.log(`Unsplash MCP server listening on http://${host}:${port}`);
    console.log(`Health check available at http://${host}:${port}/health`);
  });
}

async function main() {
  const mode = process.argv[2] || "stdio";

  switch (mode) {
    case "stdio":
      await runStdio();
      break;
    case "server":
      await runHttp(SERVER_CONFIG.host!, SERVER_CONFIG.port!);
      break;
    default:
      console.error(`Unknown mode: ${mode}\n`);
      showUsage();
      process.exit(2);
  }
}

main().catch((error) => {
  console.error("Server error:", error);
  process.exit(1);
});