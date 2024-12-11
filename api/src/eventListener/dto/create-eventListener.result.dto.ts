import { ApiProperty } from "@nestjs/swagger";
import { IsUUID } from "class-validator";
import { UUID } from "crypto";

export class CreateEventListenerResultDto {
  @ApiProperty({
    example: "221fabd2-9df0-44e7-98d4-4eda556c4143",
    description: "Id of created event listener",
  })
  @IsUUID()
  eventListenerId: UUID;
}
