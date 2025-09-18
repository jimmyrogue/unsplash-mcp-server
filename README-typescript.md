# Unsplash MCP Server (TypeScript)

A TypeScript implementation of a Model Context Protocol (MCP) server that provides a `search_photos` tool to query the Unsplash API.

## 🆕 TypeScript Migration

This is the TypeScript version of the original Go implementation, providing:

- ✅ **Type Safety**: Full TypeScript type safety with Zod runtime validation
- ✅ **Modern Development**: Better IDE support and development experience
- ✅ **API Compatibility**: 100% compatible with the original Go version
- ✅ **Enhanced Error Handling**: Structured error handling with detailed context
- ✅ **Improved Architecture**: Clean separation of concerns with modular design

## Features

- **Dual Mode Support**: Run as stdio or HTTP server
- **Parameter Validation**: Comprehensive input validation with Zod schemas
- **Error Handling**: Robust error handling with timeout management
- **Session Management**: HTTP mode supports persistent sessions
- **Health Monitoring**: Built-in health check endpoint
- **Request Logging**: Configurable request/response logging

## Project Structure

```
src/
  index.ts             # Main entry point with dual mode support
  server.ts            # MCP server instance and tool registration
  handlers/
    searchPhotos.ts    # search_photos tool handler
  schemas/
    searchPhotos.ts    # Zod validation schemas
  services/
    unsplash.ts        # Unsplash API client
  middleware/
    logging.ts         # HTTP logging middleware
  types/
    index.ts           # TypeScript type definitions
```

## Requirements

- **Node.js** 18.0.0+
- **UNSPLASH_ACCESS_KEY** environment variable

## Installation

```bash
# Install dependencies
npm install

# Build the project
npm run build

# Copy environment template
cp .env.example .env
# Edit .env and add your Unsplash API key
```

## Environment Variables

```bash
# Required: Your Unsplash API access key
UNSPLASH_ACCESS_KEY=your_unsplash_access_key_here

# Optional: Server configuration
HOST=127.0.0.1
PORT=8080
```

## Usage

### Stdio Mode (Default)

```bash
# Development
npm run stdio

# Production
npm start stdio
```

### HTTP Server Mode

```bash
# Development
npm run server

# Production
npm start server
```

### Available Scripts

```bash
npm run build        # Build TypeScript to JavaScript
npm run dev          # Run in development mode (stdio)
npm run dev server   # Run HTTP server in development
npm run dev stdio    # Run stdio mode in development
npm run start        # Run production build (stdio)
npm run test         # Run tests
npm run lint         # Run ESLint
npm run lint:fix     # Fix ESLint issues
```

## API Endpoints (HTTP Mode)

- `GET /health` - Health check endpoint
- `POST /mcp` - MCP protocol endpoint
- `GET /mcp` - Server-sent events for notifications
- `DELETE /mcp` - Session termination

### Health Check Response

```json
{
  "status": "ok",
  "service": "unsplash-mcp-server-ts",
  "version": "1.0.0",
  "timestamp": "2025-09-18T05:59:56.123Z"
}
```

## Tool: search_photos

Search for photos on Unsplash with comprehensive filtering options.

### Parameters

- `query` (string, required): Search keyword
- `page` (number, optional): Page number (default: 1)
- `per_page` (number, optional): Results per page (1-30, default: 10)
- `order_by` (string, optional): Sort order ("relevant" or "latest", default: "relevant")
- `color` (string, optional): Color filter (black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, blue)
- `orientation` (string, optional): Orientation filter (landscape, portrait, squarish)

### Response Format

```json
{
  "query": "nature",
  "page": 1,
  "per_page": 10,
  "order_by": "relevant",
  "total": 1000000,
  "total_pages": 100000,
  "results": [
    {
      "id": "photo_id",
      "description": "A beautiful nature photo",
      "alt_description": "Green forest landscape",
      "urls": {
        "raw": "https://images.unsplash.com/...",
        "full": "https://images.unsplash.com/...",
        "regular": "https://images.unsplash.com/...",
        "small": "https://images.unsplash.com/...",
        "thumb": "https://images.unsplash.com/..."
      },
      "width": 4000,
      "height": 3000
    }
  ],
  "retrieved_at": "2025-09-18T05:59:56.123Z"
}
```

## Development

### Testing

```bash
# Run all tests
npm test

# Run comprehensive server test
node test-server.js

# Manual testing
curl http://localhost:8080/health
```

### Code Quality

```bash
# Lint code
npm run lint

# Fix linting issues
npm run lint:fix

# Build and check types
npm run build
```

## Migration from Go Version

This TypeScript version maintains 100% API compatibility with the original Go implementation:

- ✅ Same tool interface and parameters
- ✅ Identical response format
- ✅ Same validation logic and error handling
- ✅ Compatible dual-mode operation
- ✅ Same environment variable configuration

You can seamlessly replace the Go version with this TypeScript implementation.

## Architecture Highlights

### Type Safety
- **Zod Schemas**: Runtime validation with compile-time type inference
- **Strict TypeScript**: Full type coverage with strict compiler settings
- **Type Definitions**: Comprehensive interfaces for all data structures

### Error Handling
- **Structured Errors**: Consistent error format with context
- **Timeout Management**: 15-second request timeout with AbortController
- **Validation Errors**: Clear parameter validation error messages

### Performance
- **Parameter Normalization**: Efficient input processing and validation
- **Session Reuse**: HTTP mode supports persistent connections
- **Minimal Dependencies**: Lean dependency tree for fast startup

## Troubleshooting

### Common Issues

1. **Missing API Key**: Ensure `UNSPLASH_ACCESS_KEY` is set in `.env`
2. **Port Conflicts**: Change `PORT` in `.env` if 8080 is occupied
3. **Build Errors**: Run `npm run build` to check TypeScript compilation
4. **Network Issues**: Verify Unsplash API accessibility

### Debug Mode

Set `NODE_ENV=development` for detailed logging:

```bash
NODE_ENV=development npm run dev server
```

## License

MIT

## Contributing

1. Fork the repository
2. Create a feature branch
3. Run tests: `npm test`
4. Submit a pull request

---

**Migration Complete!** 🎉

The TypeScript version is now ready and fully functional, providing a modern, type-safe alternative to the original Go implementation.