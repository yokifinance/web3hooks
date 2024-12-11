import { Logger, ValidationPipe, VersioningType } from "@nestjs/common";
import { NestFactory } from "@nestjs/core";
import { ConfigService } from "@nestjs/config";
import { useContainer } from "class-validator";
import helmet from "helmet";
import swaggerInit from "./swagger";
import { AppModule } from "./app.module";

export default async function (logger: Logger) {
  const app = await NestFactory.create(AppModule, {
    logger: ["log", "error", "warn", "debug"],
  });

  // Start app

  // public API must follow versioning (OPTIONAL)
  // for internal modules - remove this
  // app.enableVersioning({
  //   type: VersioningType.URI,
  //   defaultVersion: '1',
  // });
  app.useGlobalPipes(new ValidationPipe());
  app.enableShutdownHooks();
  useContainer(app.select(AppModule), { fallbackOnErrors: true });

  const configService = app.get(ConfigService);

  // web page for documentation, OPTIONAL
  await swaggerInit(app, configService);

  app.enableCors({
    origin: ["http://localhost:3000"],
    allowedHeaders: ["*"],
    methods: ["GET", "POST", "DELETE", "OPTIONS"],
    credentials: true,
  });
  // basic security
  app.use(helmet());

  const port = configService.get("application.port");
  await app.listen(port, () => {
    logger.log(`API application is up on port: ${port}`);
  });
  return app;
}
