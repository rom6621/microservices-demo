import { adClient } from "../clients.js";
import z from "zod";
import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { Code, ConnectError } from "@connectrpc/connect";
import { ErrorResponseSchema, type ErrorResponse } from "./error.js";

export const adRoutes = new OpenAPIHono();

const AdSchema = z
  .object({
    text: z.string().openapi({ example: "Tank top for sale. 20% off." }),
    redirectUri: z
      .url()
      .openapi({ example: "https://example.com/product/66VCHSJNUP" }),
  })
  .openapi("Ad");
type Ad = z.infer<typeof AdSchema>;

export const getAds = createRoute({
  method: "get",
  path: "/ads",
  tags: ["Ad"],
  request: {
    query: z.object({
      contextKeys: z
        .preprocess(
          (v) => (v === undefined ? undefined : Array.isArray(v) ? v : [v]),
          z.string().array(),
        )
        .optional(),
    }),
  },
  description: "最大2件の広告をランダムに取得する",
  responses: {
    200: {
      description: "広告 (最大2件)",
      content: {
        "application/json": {
          schema: z.object({
            ads: AdSchema.array(),
          }),
        },
      },
    },
    404: {
      description: "NotFound",
      content: {
        "application/json": {
          schema: ErrorResponseSchema,
        },
      },
    },
    500: {
      description: "Internal Server Error",
      content: {
        "application/json": {
          schema: ErrorResponseSchema,
        },
      },
    },
    503: {
      description: "Service Unavailable",
      content: {
        "application/json": {
          schema: ErrorResponseSchema,
        },
      },
    },
    504: {
      description: "Gateway Timeout",
      content: {
        "application/json": {
          schema: ErrorResponseSchema,
        },
      },
    },
  },
});

adRoutes.openapi(getAds, async (c) => {
  const q = c.req.valid("query");

  try {
    const res = await adClient.getAds({ contextKeys: q.contextKeys ?? [] });
    const ads = res.ads.map(
      (ad) =>
        ({
          text: ad.text,
          redirectUri: ad.redirectUrl,
        }) satisfies Ad,
    );

    return c.json({ ads }, 200);
  } catch (err: unknown) {
    if (err instanceof ConnectError && err.code === Code.Unavailable) {
      return c.json(
        {
          status: 503,
          type: "/ads/internal",
          title: "広告の取得に失敗しました。",
          detail:
            "サービスが一時的に利用できません。時間をおいてやり直してください。",
          instance: `/ads`,
        } satisfies ErrorResponse,
        500,
      );
    }

    if (err instanceof ConnectError && err.code === Code.DeadlineExceeded) {
      return c.json(
        {
          status: 504,
          type: "/ads/internal",
          title: "広告の取得に失敗しました。",
          detail:
            "応答に時間がかかっています。時間をおいてやり直してください。",
          instance: `/ads`,
        } satisfies ErrorResponse,
        500,
      );
    }

    return c.json(
      {
        status: 500,
        type: "/ads/internal",
        title: "広告の取得に失敗しました。",
        detail:
          "予期しないエラーが発生しました。時間をおいてやり直してください。",
        instance: `/ads`,
      } satisfies ErrorResponse,
      500,
    );
  }
});
