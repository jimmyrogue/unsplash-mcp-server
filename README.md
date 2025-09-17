# Unsplash MCP Server / Unsplash MCP 服务器

## English

This project implements a Model Context Protocol (MCP) server written in Go. The server exposes a single tool, `search_photos`, that proxies the Unsplash Search API and returns structured photo metadata. You can run the service either over standard input/output (stdio) for direct MCP integrations or as an HTTP Streamable MCP endpoint.

### Requirements
- Go 1.21 or newer
- An Unsplash access key (`UNSPLASH_ACCESS_KEY`)

### Setup
1. Export your Unsplash API key:
   ```bash
   export UNSPLASH_ACCESS_KEY=your_key_here
   ```
2. Install dependencies (handled automatically by `go run` / `go build`).

### Usage
- **Stdio mode (default):**
  ```bash
  go run .
  ```
- **HTTP mode:**
  ```bash
  go run . server -host 0.0.0.0 -port 8080
  ```
  The server listens on `http://host:port` and serves the MCP Streamable transport.

### Tool: `search_photos`
| Argument    | Type    | Description                                                       |
|-------------|---------|-------------------------------------------------------------------|
| `query`     | string  | Search keyword (required)                                         |
| `page`      | number  | Page number (default: 1)                                          |
| `per_page`  | number  | Results per page (default: 10, max: 30)                           |
| `order_by`  | string  | Sort order (`relevant` or `latest`, default: `relevant`)         |
| `color`     | string  | Color filter (optional; Unsplash-supported values)               |
| `orientation` | string | Orientation filter (`landscape`, `portrait`, `squarish`)         |

### Response
The tool returns structured JSON containing the query, pagination details, and an array of Unsplash photo objects with IDs, descriptions, dimensions, and URL variants.

---

## 中文

本项目是一个使用 Go 编写的 Model Context Protocol (MCP) 服务器。服务器暴露了一个名为 `search_photos` 的工具，用于调用 Unsplash 图片搜索 API，并返回结构化的图片元数据。服务器支持通过标准输入/输出 (stdio) 或 HTTP Streamable MCP 端点两种运行方式。

### 环境要求
- Go 1.21 或更高版本
- Unsplash 访问密钥 (`UNSPLASH_ACCESS_KEY`)

### 配置步骤
1. 设置 Unsplash API Key：
   ```bash
   export UNSPLASH_ACCESS_KEY=你的密钥
   ```
2. 依赖由 `go run` / `go build` 自动拉取。

### 使用方式
- **Stdio 模式（默认）**
  ```bash
  go run .
  ```
- **HTTP 模式**
  ```bash
  go run . server -host 0.0.0.0 -port 8080
  ```
  服务器会在 `http://host:port` 上提供 MCP Streamable 传输接口。

### 工具：`search_photos`
| 参数           | 类型   | 说明                                             |
|----------------|--------|--------------------------------------------------|
| `query`        | string | 搜索关键字（必填）                               |
| `page`         | number | 页码（默认 1）                                   |
| `per_page`     | number | 每页结果数量（默认 10，最大 30）                |
| `order_by`     | string | 排序方式（`relevant` 或 `latest`，默认 relevant）|
| `color`        | string | 颜色筛选（可选，需符合 Unsplash 支持的枚举）     |
| `orientation`  | string | 图片方向（`landscape`、`portrait`、`squarish`） |

### 返回结果
工具返回一份结构化的 JSON，其中包含查询参数、分页信息，以及一组 Unsplash 图片对象（包括 ID、描述、尺寸以及各种尺寸的图片链接）。

