import { serve } from "@hono/node-server";
import { swaggerUI } from "@hono/swagger-ui";
import { OpenAPIHono } from "@hono/zod-openapi";
import { adRoutes } from "./routes/ads.js";
import { productRoutes } from "./routes/products.js";

export const app = new OpenAPIHono();
app.get("/docs/openapi", swaggerUI({ url: "/openapi" }));
app.doc("/openapi", {
  openapi: "3.1.0",
  info: {
    version: "0.1.0",
    title: "Online Boutique BFF",
  },
});

app.get("/api/healthz", (c) => c.text("ok"));
app.route("/api", productRoutes);
app.route("/api", adRoutes);

export type AppType = typeof app;

const port = Number(process.env.PORT ?? 3000);
serve(
  {
    fetch: app.fetch,
    port,
  },
  (info) => {
    console.log(`listening on: ${info.port}`);
  },
);
