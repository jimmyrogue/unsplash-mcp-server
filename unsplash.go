package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	unsplashSearchEndpoint = "https://api.unsplash.com/search/photos"
	requestTimeout         = 15 * time.Second
	defaultPage            = 1
	defaultPerPage         = 10
	maxPerPage             = 30
	defaultOrderBy         = "relevant"
)

var (
	httpClient = &http.Client{Timeout: requestTimeout}

	allowedOrderBy = map[string]struct{}{
		"relevant": {},
		"latest":   {},
	}

	allowedColors = map[string]struct{}{
		"black_and_white": {},
		"black":           {},
		"white":           {},
		"yellow":          {},
		"orange":          {},
		"red":             {},
		"purple":          {},
		"magenta":         {},
		"green":           {},
		"teal":            {},
		"blue":            {},
	}

	allowedOrientations = map[string]struct{}{
		"landscape": {},
		"portrait":  {},
		"squarish":  {},
	}
)

// searchPhotosArgs captures the incoming tool arguments from the client.
type searchPhotosArgs struct {
	Query       string  `json:"query" jsonschema:"Search keyword"`
	Page        int     `json:"page" jsonschema:"Page number (1-based). Default: 1"`
	PerPage     int     `json:"per_page" jsonschema:"Results per page (1-30). Default: 10"`
	OrderBy     string  `json:"order_by" jsonschema:"Sort method (relevant or latest). Default: relevant"`
	Color       *string `json:"color,omitempty" jsonschema:"Color filter (black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, blue)"`
	Orientation *string `json:"orientation,omitempty" jsonschema:"Orientation filter (landscape, portrait, squarish)"`
}

type unsplashPhoto struct {
	ID             string            `json:"id"`
	Description    *string           `json:"description,omitempty"`
	AltDescription *string           `json:"alt_description,omitempty"`
	URLs           map[string]string `json:"urls"`
	Width          int               `json:"width"`
	Height         int               `json:"height"`
}

type unsplashSearchAPIResponse struct {
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
	Results    []unsplashPhoto `json:"results"`
}

// searchPhotosResult is returned as structured tool output for the client.
type searchPhotosResult struct {
	Query       string          `json:"query"`
	Page        int             `json:"page"`
	PerPage     int             `json:"per_page"`
	OrderBy     string          `json:"order_by"`
	Color       string          `json:"color,omitempty"`
	Orientation string          `json:"orientation,omitempty"`
	Total       int             `json:"total"`
	TotalPages  int             `json:"total_pages"`
	Results     []unsplashPhoto `json:"results"`
	RetrievedAt time.Time       `json:"retrieved_at"`
}

// normalizeEnum trims, lowercases, and validates a required enum value.
func normalizeEnum(value, field, fallback string, allowed map[string]struct{}) (string, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		v = fallback
	}
	if _, ok := allowed[v]; !ok {
		return "", fmt.Errorf("invalid %s value: %s", field, v)
	}
	return v, nil
}

// normalizeOptionalEnum validates an optional enum value, treating blank input as unset.
func normalizeOptionalEnum(value *string, field string, allowed map[string]struct{}) (string, error) {
	if value == nil {
		return "", nil
	}
	v := strings.ToLower(strings.TrimSpace(*value))
	if v == "" {
		return "", nil
	}
	if _, ok := allowed[v]; !ok {
		return "", fmt.Errorf("invalid %s value: %s", field, v)
	}
	return v, nil
}

// clampPerPage keeps the per_page argument within Unsplash's supported bounds.
func clampPerPage(value int) int {
	switch {
	case value <= 0:
		return defaultPerPage
	case value > maxPerPage:
		return maxPerPage
	default:
		return value
	}
}

// handleSearchPhotos calls the Unsplash Search API and returns structured results.
func handleSearchPhotos(ctx context.Context, _ *mcp.CallToolRequest, args searchPhotosArgs) (*mcp.CallToolResult, searchPhotosResult, error) {
	query := strings.TrimSpace(args.Query)
	if query == "" {
		return nil, searchPhotosResult{}, fmt.Errorf("query is required")
	}

	page := args.Page
	if page < 1 {
		page = defaultPage
	}

	perPage := clampPerPage(args.PerPage)

	orderBy, err := normalizeEnum(args.OrderBy, "order_by", defaultOrderBy, allowedOrderBy)
	if err != nil {
		return nil, searchPhotosResult{}, err
	}

	color, err := normalizeOptionalEnum(args.Color, "color", allowedColors)
	if err != nil {
		return nil, searchPhotosResult{}, err
	}

	orientation, err := normalizeOptionalEnum(args.Orientation, "orientation", allowedOrientations)
	if err != nil {
		return nil, searchPhotosResult{}, err
	}

	accessKey := strings.TrimSpace(os.Getenv("UNSPLASH_ACCESS_KEY"))
	if accessKey == "" {
		return nil, searchPhotosResult{}, fmt.Errorf("missing UNSPLASH_ACCESS_KEY environment variable")
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	params.Set("order_by", orderBy)
	if color != "" {
		params.Set("color", color)
	}
	if orientation != "" {
		params.Set("orientation", orientation)
	}

	requestCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(requestCtx, http.MethodGet, unsplashSearchEndpoint, nil)
	if err != nil {
		return nil, searchPhotosResult{}, fmt.Errorf("unable to create Unsplash request: %w", err)
	}

	httpReq.Header.Set("Accept-Version", "v1")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", accessKey))
	httpReq.URL.RawQuery = params.Encode()

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, searchPhotosResult{}, fmt.Errorf("request to Unsplash timed out")
		}
		return nil, searchPhotosResult{}, fmt.Errorf("request to Unsplash failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, searchPhotosResult{}, fmt.Errorf("failed to read Unsplash response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = resp.Status
		}
		if len(message) > 512 {
			message = message[:512] + "..."
		}
		return nil, searchPhotosResult{}, fmt.Errorf("Unsplash API error (%d): %s", resp.StatusCode, message)
	}

	var apiResp unsplashSearchAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, searchPhotosResult{}, fmt.Errorf("failed to decode Unsplash response: %w", err)
	}

	result := searchPhotosResult{
		Query:       query,
		Page:        page,
		PerPage:     perPage,
		OrderBy:     orderBy,
		Total:       apiResp.Total,
		TotalPages:  apiResp.TotalPages,
		Results:     apiResp.Results,
		RetrievedAt: time.Now().UTC(),
	}

	if color != "" {
		result.Color = color
	}
	if orientation != "" {
		result.Orientation = orientation
	}

	return nil, result, nil
}
