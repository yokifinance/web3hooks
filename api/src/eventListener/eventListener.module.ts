import { Module } from "@nestjs/common";
import { EventListenerService } from "./eventListener.service";
import { EventListenerController } from "./eventListener.controller";

@Module({
  controllers: [EventListenerController],
  providers: [EventListenerService],
})
export class EventListenerModule {}
