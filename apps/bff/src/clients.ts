import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import { AdService, ProductCatalogService } from "./gen/demo_pb.js";

const productCatalogTransport = createGrpcTransport({
  baseUrl:
    process.env.PRODUCT_CATALOG_SERVICE_ENDPOINT ?? "http://localhost:3550",
  defaultTimeoutMs: 3000,
});

export const productCatalogClient = createClient(
  ProductCatalogService,
  productCatalogTransport,
);

const adTransport = createGrpcTransport({
  baseUrl: process.env.AD_SERVICE_ENDPOINT ?? "http://localhost:9555",
  defaultTimeoutMs: 3000,
});

export const adClient = createClient(AdService, adTransport);
