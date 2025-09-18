# Unsplash MCP Server (TypeScript)

> 🚨 **This project has been migrated to TypeScript!**
>
> The original Go implementation has been replaced with a modern TypeScript version. See [README-typescript.md](./README-typescript.md) for complete documentation.

A TypeScript implementation of a Model Context Protocol (MCP) server that provides a `search_photos` tool to query the Unsplash API.

## Quick Start

```bash
# Install dependencies
npm install

# Set up environment
cp .env.example .env
# Edit .env and add your UNSPLASH_ACCESS_KEY

# Run in development mode
npm run dev server    # HTTP mode
npm run dev stdio     # stdio mode

# Build for production
npm run build
npm start
```

## Features

- ✅ **Type Safety**: Full TypeScript with Zod validation
- ✅ **Dual Mode**: stdio and HTTP server support
- ✅ **API Compatible**: 100% compatible with original Go version
- ✅ **Modern Stack**: Express.js, TypeScript, ESM modules
- ✅ **Health Monitoring**: Built-in health endpoints
- ✅ **Session Management**: Persistent HTTP sessions

## Documentation

For complete documentation, installation instructions, and API reference, see:

**👉 [README-typescript.md](./README-typescript.md)**

## Migration Notes

The TypeScript version provides identical functionality to the original Go implementation:

- Same `search_photos` tool interface
- Same parameter validation and error handling
- Same dual-mode operation (stdio/HTTP)
- Same environment variable configuration
- Same response format

You can seamlessly replace any Go version usage with this TypeScript implementation.

---

**🎉 TypeScript Migration Complete!**

The project is now running on modern TypeScript with enhanced type safety, better development experience, and improved maintainability.