import { Column, Entity, PrimaryColumn } from "typeorm";

@Entity("supported_chains")
export class SupportedChainEntity {
  @PrimaryColumn()
  chain: number;

  @Column({ name: "name" })
  name: string;

  @Column({ name: "rpc_url" })
  rpcUrl: string;
}
