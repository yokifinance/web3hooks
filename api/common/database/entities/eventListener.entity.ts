// clients.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  ManyToOne,
  JoinColumn,
  Index,
} from "typeorm";
import { SupportedChainEntity } from "./supportedChain.entity";
import { ClientEntity } from "./client.entity";

@Entity({ name: "event_listeners" })
export class EventListenerEntity {
  @PrimaryGeneratedColumn("uuid", { name: "id" })
  id: string;

  @Column({ name: "chain" })
  @Index("event_listeners_chain_idx")
  chain: number;

  @ManyToOne(() => SupportedChainEntity)
  @JoinColumn({
    name: "chain",
    foreignKeyConstraintName: "event_listeners_chain_id_fk",
  })
  chainEntity: SupportedChainEntity;

  @Column({ name: "client_id" })
  @Index("event_listeners_client_id_idx")
  clientId: string;

  @ManyToOne(() => ClientEntity)
  @JoinColumn({
    name: "client_id",
    foreignKeyConstraintName: "event_listeners_client_id_fk",
  })
  client: ClientEntity;

  @Column("text", { name: "address" })
  address: string;

  @Column("text", { name: "webhook_url" })
  webhookUrl: string;

  @Column({ name: "created_timestamp", type: "timestamp with time zone" })
  createdTimestamp: Date;

  @Column({})
  @Index("event_listeners_active_idx")
  active: boolean;
}
