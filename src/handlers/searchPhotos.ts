import { UnsplashService } from "../services/unsplash.js";
import { searchPhotosInputSchema, SearchPhotosInput } from "../schemas/searchPhotos.js";

export class SearchPhotosHandler {
  private unsplashService: UnsplashService;

  constructor(unsplashService: UnsplashService) {
    this.unsplashService = unsplashService;
  }

  async handle(input: SearchPhotosInput) {
    try {
      // 验证和规范化输入参数
      const validatedArgs = searchPhotosInputSchema.parse(input);

      // 调用 Unsplash API
      const result = await this.unsplashService.searchPhotos(validatedArgs);

      return {
        content: [{
          type: "text" as const,
          text: JSON.stringify(result, null, 2)
        }]
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Unknown error occurred";

      return {
        content: [{
          type: "text" as const,
          text: `Error: ${errorMessage}`
        }],
        isError: true
      };
    }
  }
}