package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// 简单的HTTP包装器，避免MCP SDK的问题
func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Simple HTTP MCP server listening on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
MCP Unsplash Server

Endpoints:
- GET /health - Health check
- GET /search?query=cats&page=1&per_page=10 - Search photos

Example:
curl "http://120.27.151.229:9999/search?query=nature&per_page=5"
`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "unsplash-mcp-server",
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPageStr := r.URL.Query().Get("per_page")
	perPage := 10
	if perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 30 {
			perPage = pp
		}
	}

	orderBy := r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = "relevant"
	}

	color := r.URL.Query().Get("color")
	orientation := r.URL.Query().Get("orientation")

	// 构建searchPhotosArgs
	args := searchPhotosArgs{
		Query:   query,
		Page:    page,
		PerPage: perPage,
		OrderBy: orderBy,
	}
	if color != "" {
		args.Color = &color
	}
	if orientation != "" {
		args.Orientation = &orientation
	}

	// 调用原来的搜索函数
	_, result, err := handleSearchPhotos(context.Background(), nil, args)
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回JSON结果
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}