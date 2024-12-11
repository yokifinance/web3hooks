import { DocumentBuilder, SwaggerModule } from "@nestjs/swagger";
import { INestApplication } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { version } from "../package.json";

export default async function swaggerInit(
  app: INestApplication,
  config: ConfigService
) {
  const documentBuild = new DocumentBuilder()
    .setTitle("Yoki-API")
    .setDescription("All methods available on Yoki-backend api")
    .setVersion(`v${version}`)
    .addGlobalParameters({ name: "secret-key", in: "header", required: true })
    .build();

  const document = SwaggerModule.createDocument(app, documentBuild, {
    deepScanRoutes: true,
    extraModels: [],
  });

  SwaggerModule.setup("swagger", app, document, {
    explorer: true,
    customSiteTitle: config.get("market.title"),
  });
}
