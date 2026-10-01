import { ConnectError, Code } from "@connectrpc/connect";
import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import z from "zod";
import { productCatalogClient } from "../clients.js";
import { ErrorResponseSchema, type ErrorResponse } from "./error.js";

export const productRoutes = new OpenAPIHono();

const MoneySchema = z
  .object({
    currencyCode: z.string().openapi({ example: "USD" }),
    units: z.string().openapi({ example: "19" }),
    nanos: z.number().openapi({ example: 990000000 }),
  })
  .openapi("Money");
type Money = z.infer<typeof MoneySchema>;

const ProductSchema = z
  .object({
    id: z.string().openapi({ example: "OLJCESPC7Z" }),
    name: z.string().openapi({ example: "Sunglasses" }),
    description: z.string().openapi({
      example:
        "Add a modern touch to your outfits with these sleek aviator sunglasses.",
    }),
    picture: z.url().openapi({
      example: "https://example.com/static/img/products/sunglasses.jpg",
    }),
    price: MoneySchema,
    categories: z
      .string()
      .array()
      .openapi({ example: ["accessories"] }),
  })
  .openapi("Product");
type Product = z.infer<typeof ProductSchema>;

const getProductRoute = createRoute({
  method: "get",
  path: "/products/{id}",
  tags: ["Product"],
  description: "商品を取得する",
  request: {
    params: z.object({
      id: z.string(),
    }),
  },
  responses: {
    200: {
      description: "商品",
      content: {
        "application/json": {
          schema: ProductSchema,
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

productRoutes.openapi(getProductRoute, async (c) => {
  const req = c.req.param();

  try {
    const res = await productCatalogClient.getProduct({
      id: req.id,
    });

    const price = res.priceUsd;
    if (!price) {
      return c.json(
        {
          type: "abount:blank",
          status: 500,
          title: "商品の取得に失敗しました。",
          detail: "商品の価格が設定されていません。",
          instance: `/products/${req.id}`,
        } satisfies ErrorResponse,
        500,
      );
    }

    return c.json(
      {
        id: res.id,
        name: res.name,
        description: res.description,
        picture: res.picture,
        price: {
          currencyCode: price.currencyCode,
          units: price.units.toString(),
          nanos: price.nanos,
        },
        categories: res.categories,
      } satisfies Product,
      200,
    );
  } catch (err: unknown) {
    if (err instanceof ConnectError && err.code === Code.NotFound) {
      return c.json(
        {
          status: 404,
          type: "/products/not-found",
          title: "商品の取得に失敗しました",
          detail: "指定された商品は存在しないか、削除された可能性があります。",
          instance: `/products/${req.id}`,
        } satisfies ErrorResponse,
        404,
      );
    }

    if (err instanceof ConnectError && err.code === Code.Unavailable) {
      return c.json(
        {
          status: 503,
          type: "/products/internal",
          title: "商品の取得に失敗しました。",
          detail:
            "サービスが一時的に利用できません。時間をおいてやり直してください。",
          instance: `/products/${req.id}`,
        } satisfies ErrorResponse,
        500,
      );
    }

    if (err instanceof ConnectError && err.code === Code.DeadlineExceeded) {
      return c.json(
        {
          status: 504,
          type: "/products/internal",
          title: "商品の取得に失敗しました。",
          detail:
            "応答に時間がかかっています。時間をおいてやり直してください。",
          instance: `/products/${req.id}`,
        } satisfies ErrorResponse,
        500,
      );
    }

    return c.json(
      {
        status: 500,
        type: "/products/internal",
        title: "商品の取得に失敗しました。",
        detail:
          "予期しないエラーが発生しました。時間をおいてやり直してください。",
        instance: `/products/${req.id}`,
      } satisfies ErrorResponse,
      500,
    );
  }
});

export const listProducts = createRoute({
  method: "get",
  path: "/products",
  tags: ["Product"],
  description: "商品一覧を取得する",
  responses: {
    200: {
      description: "商品一覧",
      content: {
        "application/json": {
          schema: z.object({
            products: ProductSchema.array(),
          }),
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

class NoPriceError extends Error {
  _name: string;

  constructor() {
    super();
    this._name = "no_price_error";
  }
}

productRoutes.openapi(listProducts, async (c) => {
  try {
    const res = await productCatalogClient.listProducts({});

    const products = res.products.map((p) => {
      const price = p.priceUsd;

      if (!price) {
        throw new NoPriceError();
      }

      return {
        id: p.id,
        name: p.name,
        description: p.description,
        picture: p.picture,
        price: {
          currencyCode: price.currencyCode,
          units: price.units.toString(),
          nanos: price.nanos,
        } satisfies Money,
        categories: p.categories,
      } satisfies Product;
    });

    return c.json({ products }, 200);
  } catch (err: unknown) {
    if (err instanceof NoPriceError) {
      return c.json(
        {
          type: "abount:blank",
          status: 500,
          title: "商品の取得に失敗しました。",
          detail: "商品の価格が設定されていません。",
          instance: "/products",
        } satisfies ErrorResponse,
        500,
      );
    }

    if (err instanceof ConnectError && err.code === Code.Unavailable) {
      return c.json(
        {
          status: 503,
          type: "/products/internal",
          title: "商品の取得に失敗しました。",
          detail:
            "サービスが一時的に利用できません。時間をおいてやり直してください。",
          instance: "/products",
        } satisfies ErrorResponse,
        500,
      );
    }

    if (err instanceof ConnectError && err.code === Code.DeadlineExceeded) {
      return c.json(
        {
          status: 504,
          type: "/products/internal",
          title: "商品の取得に失敗しました。",
          detail:
            "応答に時間がかかっています。時間をおいてやり直してください。",
          instance: "/products",
        } satisfies ErrorResponse,
        500,
      );
    }

    return c.json(
      {
        status: 500,
        type: "/products/internal",
        title: "商品の取得に失敗しました。",
        detail:
          "予期しないエラーが発生しました。時間をおいてやり直してください。",
        instance: "/products",
      } satisfies ErrorResponse,
      500,
    );
  }
});
