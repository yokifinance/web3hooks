import { ApiProperty } from "@nestjs/swagger";
import { Transform } from "class-transformer";
import { IsInt, IsNumber, IsOptional, IsUrl, Max, Min } from "class-validator";

export class CreateWebhookDto {
  static DefaultErrorCount = 18;

  @IsUrl({ protocols: ["http", "https"], require_tld: false })
  @ApiProperty({
    example: "http://localhost/webhookhandler",
    description: "Webhook to call",
    required: true,
  })
  webhookUrl: string;

  @ApiProperty({
    example: { x: 10, text: "hello" },
    description: "JSON object to pass to webhook body",
    required: false,
    default: "",
  })
  @IsOptional()
  webhookBody: Object;

  @ApiProperty({
    example: "18",
    description:
      "Number of errors the service can tolerate when calling webhook (not 20x error code). In production, 18 is recommended (about 12 hours). Max is 30 (about 24 hours)",
    required: false,
    default: CreateWebhookDto.DefaultErrorCount,
  })
  @IsOptional()
  @IsInt()
  @Min(1)
  @Max(30)
  maxErrorCount: number;

  @IsUrl({ protocols: ["http", "https"], require_tld: false })
  @ApiProperty({
    example: "http://localhost/statusurl",
    description:
      "Webhook URL to call to signal job's result (if required). In body, object { webhookJobId: string, success: boolean, errorMessage?: string } is passed",
    required: false,
    nullable: true,
  })
  @IsOptional()
  resultWebhookUrl?: string;
}
