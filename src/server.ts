import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import { ServerConfig } from "./types/index.js";
import { UnsplashService } from "./services/unsplash.js";
import { SearchPhotosHandler } from "./handlers/searchPhotos.js";

export function createMcpServer(config: ServerConfig): McpServer {
  const server = new McpServer({
    name: config.name,
    version: config.version,
  });

  // 注册 search_photos 工具
  registerSearchPhotosTool(server);

  return server;
}

function registerSearchPhotosTool(server: McpServer) {
  const accessKey = process.env.UNSPLASH_ACCESS_KEY;
  if (!accessKey) {
    console.warn("UNSPLASH_ACCESS_KEY not found, search_photos tool will not be available");
    return;
  }

  const unsplashService = new UnsplashService(accessKey);
  const handler = new SearchPhotosHandler(unsplashService);

  server.registerTool(
    "search_photos",
    {
      title: "Search Unsplash Photos",
      description: "Search for photos on Unsplash using various filters",
      inputSchema: {
        query: z.string().describe("Search keyword"),
        page: z.number().optional().describe("Page number (1-based). Default: 1"),
        per_page: z.number().optional().describe("Results per page (1-30). Default: 10"),
        order_by: z.enum(["relevant", "latest"]).optional().describe("Sort method. Default: relevant"),
        color: z.enum([
          "black_and_white", "black", "white", "yellow", "orange", "red",
          "purple", "magenta", "green", "teal", "blue"
        ]).optional().describe("Color filter"),
        orientation: z.enum(["landscape", "portrait", "squarish"]).optional().describe("Orientation filter")
      }
    },
    async (args) => {
      return await handler.handle(args as any);
    }
  );
}
