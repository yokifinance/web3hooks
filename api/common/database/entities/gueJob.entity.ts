import { Entity, PrimaryColumn, Column, Index } from "typeorm";

@Entity({ name: "gue_jobs" })
@Index("gue_jobs_selector_idx", ["queue", "runAt", "priority"])
export class GueJobEntity {
  @PrimaryColumn({ name: "job_id", type: "text" })
  jobId: string;

  @Column({ name: "priority", type: "smallint" })
  priority: number;

  @Column({ name: "run_at", type: "timestamp with time zone" })
  runAt: Date;

  @Column({ name: "job_type", type: "text" })
  jobType: string;

  @Column({ name: "args", type: "bytea" })
  args: Buffer | null;

  @Column({ name: "error_count", type: "integer", default: 0 })
  errorCount: number;

  @Column({ name: "last_error", type: "text", nullable: true })
  lastError: string | null;

  @Column({ name: "queue", type: "text" })
  queue: string;

  @Column("timestamp with time zone", { name: "created_at" })
  createdAt: Date;

  @Column("timestamp with time zone", { name: "updated_at" })
  updatedAt: Date;
}
