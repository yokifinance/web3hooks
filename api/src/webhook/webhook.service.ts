import {
  BadRequestException,
  Injectable,
  UnauthorizedException,
} from "@nestjs/common";
import { CreateWebhookDto } from "./dto/create-webhook.dto";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { GueJobEntity } from "common/database/entities/gueJob.entity";
import { UUID } from "crypto";
import { CreateWebhookResultDto } from "./dto/create-webhook.result.dto";
import { ulid } from "ulid";
import { buffer } from "stream/consumers";

@Injectable()
export class WebhookService {
  constructor(
    @InjectRepository(GueJobEntity)
    private gueJobRepository: Repository<GueJobEntity>
  ) {}

  async create(createWebhookDto: CreateWebhookDto, _: UUID) {
    class JobArgsWrapper {
      WebhookUrl: string;
      MaxErrorCount: number;
      Args: string;
      ResultWebhookUrl?: string;
    }

    const buf = Buffer.from(JSON.stringify(createWebhookDto.webhookBody ?? {}));
    const wrapper: JobArgsWrapper = {
      WebhookUrl: createWebhookDto.webhookUrl,
      Args: buf.toString("base64"), // encode as base64
      MaxErrorCount:
        createWebhookDto.maxErrorCount ?? CreateWebhookDto.DefaultErrorCount,
      ResultWebhookUrl: createWebhookDto.resultWebhookUrl,
    };
    const wrappedArgs = JSON.stringify(wrapper);

    const date = new Date();
    const job = this.gueJobRepository.create({
      jobId: ulid(),
      createdAt: date,
      updatedAt: date,
      runAt: date,
      queue: "webhookService_queue",
      jobType: "webhookService_jobType",
      priority: 0,
      args: Buffer.from(wrappedArgs, "utf8"),
    });

    await this.gueJobRepository.save(job);

    return {
      webhookJobId: job.jobId,
    } as CreateWebhookResultDto;
  }

  async checkWebhookJobExists(webhookJobId: string) {
    return this.gueJobRepository.exist({ where: { jobId: webhookJobId } });
  }
}
