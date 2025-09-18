import fetch from "node-fetch";
import { UnsplashSearchResponse, SearchPhotosResult } from "../types/index.js";
import { SearchPhotosArgs } from "../schemas/searchPhotos.js";

const UNSPLASH_SEARCH_ENDPOINT = "https://api.unsplash.com/search/photos";
const REQUEST_TIMEOUT = 15000; // 15 seconds

export class UnsplashService {
  private accessKey: string;

  constructor(accessKey: string) {
    if (!accessKey || accessKey.trim() === "") {
      throw new Error("UNSPLASH_ACCESS_KEY environment variable is required");
    }
    this.accessKey = accessKey.trim();
  }

  async searchPhotos(args: SearchPhotosArgs): Promise<SearchPhotosResult> {
    const params = new URLSearchParams();
    params.set("query", args.query);
    params.set("page", args.page.toString());
    params.set("per_page", args.per_page.toString());
    params.set("order_by", args.order_by);

    if (args.color) {
      params.set("color", args.color);
    }
    if (args.orientation) {
      params.set("orientation", args.orientation);
    }

    const url = `${UNSPLASH_SEARCH_ENDPOINT}?${params.toString()}`;

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT);

      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Accept-Version": "v1",
          "Authorization": `Client-ID ${this.accessKey}`,
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const errorText = await response.text().catch(() => response.statusText);
        const message = errorText.length > 512 ? errorText.slice(0, 512) + "..." : errorText;
        throw new Error(`Unsplash API error (${response.status}): ${message}`);
      }

      const data = await response.json() as UnsplashSearchResponse;

      const result: SearchPhotosResult = {
        query: args.query,
        page: args.page,
        per_page: args.per_page,
        order_by: args.order_by,
        total: data.total,
        total_pages: data.total_pages,
        results: data.results,
        retrieved_at: new Date().toISOString(),
      };

      if (args.color) {
        result.color = args.color;
      }
      if (args.orientation) {
        result.orientation = args.orientation;
      }

      return result;
    } catch (error) {
      if (error instanceof Error) {
        if (error.name === "AbortError") {
          throw new Error("Request to Unsplash timed out");
        }
        if (error.message.includes("fetch")) {
          throw new Error(`Request to Unsplash failed: ${error.message}`);
        }
        throw error;
      }
      throw new Error("Unknown error occurred while fetching from Unsplash");
    }
  }
}