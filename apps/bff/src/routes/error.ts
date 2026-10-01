import z from "zod";

export const ErrorResponseSchema = z.object({
  type: z.string(),
  status: z.number().int(),
  title: z.string(),
  detail: z.string(),
  instance: z.string(),
});
export type ErrorResponse = z.infer<typeof ErrorResponseSchema>;
