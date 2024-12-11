import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Headers,
  Param,
  Delete,
  HttpStatus,
} from "@nestjs/common";
import { EventListenerService } from "./eventListener.service";
import { ApiOperation, ApiResponse } from "@nestjs/swagger";
import { CreateEventListenerResultDto } from "./dto/create-eventListener.result.dto";
import { CreateEventListenerDto } from "./dto/create-eventListener.dto";
import { GetEventListenerResultDto } from "./dto/get-eventListener.result.dto";
import { UUID } from "crypto";
import { AuthClient } from "common/auth/dto/auth.client";
import { Client } from "common/auth/client.decorator";
import { StopEventListenerDto } from "./dto/stop-eventListener.dto";

@Controller("eventListener")
export class EventListenerController {
  constructor(private readonly eventListenerService: EventListenerService) {}

  @Post()
  @ApiResponse({ type: CreateEventListenerResultDto })
  async create(
    @Body() createEventListenerDto: CreateEventListenerDto,
    @Client() client: AuthClient
  ) {
    return await this.eventListenerService.create(
      createEventListenerDto,
      client.id
    );
  }

  @Post("/stop")
  @ApiResponse({})
  async stop(
    @Body() stopEventListenerDto: StopEventListenerDto,
    @Client() client: AuthClient
  ) {
    return await this.eventListenerService.stop(
      stopEventListenerDto.eventListenerId,
      client.id
    );
  }

  @Get()
  @ApiOperation({
    summary: "Get all event listeners for current client",
  })
  @ApiResponse({
    status: HttpStatus.OK,
    type: GetEventListenerResultDto,
    isArray: true,
  })
  findAll(@Client() client: AuthClient) {
    return this.eventListenerService.findAll(client.id);
  }

  @Get("/:eventListenerId")
  @ApiOperation({
    summary: "Get event listener by id",
  })
  @ApiResponse({ status: HttpStatus.OK, type: GetEventListenerResultDto })
  findOne(
    @Param("eventListenerId") eventListenerId: UUID,
    @Client() client: AuthClient
  ) {
    return this.eventListenerService.findOne(eventListenerId, client.id);
  }
}
