import { z } from "zod";

// 定义允许的枚举值
export const allowedOrderBy = ["relevant", "latest"] as const;
export const allowedColors = [
  "black_and_white",
  "black",
  "white",
  "yellow",
  "orange",
  "red",
  "purple",
  "magenta",
  "green",
  "teal",
  "blue"
] as const;
export const allowedOrientations = ["landscape", "portrait", "squarish"] as const;

// 常量
export const DEFAULT_PAGE = 1;
export const DEFAULT_PER_PAGE = 10;
export const MAX_PER_PAGE = 30;
export const DEFAULT_ORDER_BY = "relevant";

// 参数验证 Schema
export const searchPhotosSchema = z.object({
  query: z.string()
    .min(1, "Query is required")
    .max(100, "Query must be 100 characters or less"),
  page: z.number()
    .int()
    .min(1, "Page must be at least 1")
    .optional()
    .default(DEFAULT_PAGE),
  per_page: z.number()
    .int()
    .min(1, "Per page must be at least 1")
    .max(MAX_PER_PAGE, `Per page must be at most ${MAX_PER_PAGE}`)
    .optional()
    .default(DEFAULT_PER_PAGE),
  order_by: z.enum(allowedOrderBy)
    .optional()
    .default(DEFAULT_ORDER_BY),
  color: z.enum(allowedColors)
    .optional(),
  orientation: z.enum(allowedOrientations)
    .optional()
});

// 输入数据预处理和规范化
export const searchPhotosInputSchema = z.object({
  query: z.string().transform(val => val.trim()),
  page: z.union([z.number(), z.string()]).transform(val => {
    const num = typeof val === "string" ? parseInt(val, 10) : val;
    return isNaN(num) || num < 1 ? DEFAULT_PAGE : num;
  }),
  per_page: z.union([z.number(), z.string()]).transform(val => {
    const num = typeof val === "string" ? parseInt(val, 10) : val;
    if (isNaN(num) || num < 1) return DEFAULT_PER_PAGE;
    return Math.min(num, MAX_PER_PAGE);
  }),
  order_by: z.string().optional().transform(val => {
    if (!val) return DEFAULT_ORDER_BY;
    const normalized = val.toLowerCase().trim();
    return allowedOrderBy.includes(normalized as any) ? normalized : DEFAULT_ORDER_BY;
  }),
  color: z.string().optional().transform(val => {
    if (!val) return undefined;
    const normalized = val.toLowerCase().trim();
    return allowedColors.includes(normalized as any) ? normalized : undefined;
  }),
  orientation: z.string().optional().transform(val => {
    if (!val) return undefined;
    const normalized = val.toLowerCase().trim();
    return allowedOrientations.includes(normalized as any) ? normalized : undefined;
  })
}).transform(data => searchPhotosSchema.parse(data));

export type SearchPhotosInput = z.input<typeof searchPhotosInputSchema>;
export type SearchPhotosArgs = z.output<typeof searchPhotosInputSchema>;