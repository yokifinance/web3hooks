import { ApiProperty } from "@nestjs/swagger";
import { IsEthereumAddress, IsNumber, IsUrl } from "class-validator";
import { IsUUID } from "class-validator";
import { UUID } from "crypto";

export class GetEventListenerResultDto {
  @ApiProperty({
    example: "221fabd2-9df0-44e7-98d4-4eda556c4143",
    description: "Id of event listener",
  })
  @IsUUID()
  id: UUID;

  @IsNumber()
  @ApiProperty({
    example: 137,
  })
  chain: number;

  @IsEthereumAddress()
  @ApiProperty({
    example: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    description: "address to listen for events",
  })
  address: string;

  @IsUrl({ protocols: ["http", "https"], require_tld: false })
  @ApiProperty({
    example: "http://localhost/wehhookhandler",
    description: "Webhook to call back",
  })
  webhookUrl: string;

  @ApiProperty({
    example: true,
    description: "Whether the listener is active",
  })
  active: boolean;
}
