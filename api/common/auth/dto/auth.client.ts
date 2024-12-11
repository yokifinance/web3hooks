import { ApiProperty } from "@nestjs/swagger";
import { UUID } from "crypto";

export class AuthClient {
  @ApiProperty()
  id: UUID;

  @ApiProperty()
  name: string;
}
