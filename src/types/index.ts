export interface UnsplashPhoto {
  id: string;
  description?: string;
  alt_description?: string;
  urls: Record<string, string>;
  width: number;
  height: number;
}

export interface UnsplashSearchResponse {
  total: number;
  total_pages: number;
  results: UnsplashPhoto[];
}

export interface SearchPhotosArgs {
  query: string;
  page?: number;
  per_page?: number;
  order_by?: string;
  color?: string;
  orientation?: string;
}

export interface SearchPhotosResult {
  query: string;
  page: number;
  per_page: number;
  order_by: string;
  color?: string;
  orientation?: string;
  total: number;
  total_pages: number;
  results: UnsplashPhoto[];
  retrieved_at: string;
}

export interface ServerConfig {
  name: string;
  version: string;
  host?: string;
  port?: number;
}