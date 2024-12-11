import { Module } from "@nestjs/common";
import { getGlobalConfigModule, DatabaseModule } from "../common";
import { EventListenerModule } from "./eventListener/eventListener.module";
import { APP_GUARD } from "@nestjs/core";
import { AuthGuard } from "common/auth/auth.guard";
import { WebhookModule } from "./webhook/webhook.module";

@Module({
  imports: [
    getGlobalConfigModule({}),
    DatabaseModule.forRoot(),
    EventListenerModule,
    WebhookModule,
  ],
  providers: [{ provide: APP_GUARD, useClass: AuthGuard }],
})
export class AppModule {}
