import { ApiProperty } from "@nestjs/swagger";

export class CreateWebhookResultDto {
  @ApiProperty({
    example: "01HJNGSXF8N96V5S7MGXBT88N1",
    description:
      "Id of created webhook job id (can be used for tracking purposes)",
  })
  webhookJobId: string;
}
