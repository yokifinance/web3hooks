import { ApiProperty } from "@nestjs/swagger";
import { IsEthereumAddress, IsNumber, IsUrl } from "class-validator";

export class CreateEventListenerDto {
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
}
