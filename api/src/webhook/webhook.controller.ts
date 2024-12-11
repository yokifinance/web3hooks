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
import { WebhookService } from "./webhook.service";
import { CreateWebhookDto } from "./dto/create-webhook.dto";
import { ApiOperation, ApiResponse } from "@nestjs/swagger";
import { CreateWebhookResultDto } from "./dto/create-webhook.result.dto";
import { Client } from "common/auth/client.decorator";
import { AuthClient } from "common/auth/dto/auth.client";

@Controller("webhook")
export class WebhookController {
  constructor(private readonly webhookService: WebhookService) {}

  @Post()
  @ApiResponse({ type: CreateWebhookResultDto })
  async create(
    @Body() createWebhookDto: CreateWebhookDto,
    @Client() client: AuthClient
  ) {
    return await this.webhookService.create(createWebhookDto, client.id);
  }

  @Get("/:webhookJobId")
  @ApiOperation({
    summary:
      "Check if webhook job exists in the database. If exists a job is not finished yet nor expired",
  })
  @ApiResponse({ status: HttpStatus.OK })
  async findOne(@Param("webhookJobId") webhookJobId: string) {
    return await this.webhookService.checkWebhookJobExists(webhookJobId);
  }
}
