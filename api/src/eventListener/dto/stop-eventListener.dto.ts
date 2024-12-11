import { ApiProperty } from "@nestjs/swagger";
import { IsEthereumAddress, IsNumber, IsUUID, IsUrl } from "class-validator";
import { UUID } from "crypto";

export class StopEventListenerDto {
  @ApiProperty({
    example: "221fabd2-9df0-44e7-98d4-4eda556c4143",
    description: "Id of listener to be stopped",
  })
  @IsUUID()
  eventListenerId: UUID;
}
