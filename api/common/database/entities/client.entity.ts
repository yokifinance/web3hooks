// clients.entity.ts
import {
  Entity,
  Column,
  PrimaryColumn,
  PrimaryGeneratedColumn,
  Unique,
} from "typeorm";

@Entity({ name: "clients" })
export class ClientEntity {
  @PrimaryGeneratedColumn("uuid", { name: "id" })
  id: string;

  @Column({ name: "name" })
  name: string;

  @Column({ name: "secret_key", unique: true })
  secretKey: string;
}
