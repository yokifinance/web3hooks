// clients.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  ManyToOne,
  JoinColumn,
  Index,
  Unique,
  PrimaryColumn,
  OneToOne,
} from "typeorm";
import { SupportedChainEntity } from "./supportedChain.entity";

@Entity({ name: "last_processed_blocks" })
export class LastProcessedBlockEntity {
  @PrimaryColumn()
  chain: number;

  @OneToOne(() => SupportedChainEntity)
  @JoinColumn({
    name: "chain",
    foreignKeyConstraintName: "last_events_chain_id_fk",
  })
  chainEntity: SupportedChainEntity;

  @Column("numeric", { name: "last_processed_block_number" })
  lastProcessedBlockNumber: BigInt;
}
