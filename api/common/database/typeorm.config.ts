import { DataSource, DataSourceOptions } from "typeorm";
import * as dotenv from "dotenv";
import * as path from "path";

dotenv.config({ path: path.join(__dirname, "..", "..", "..", ".env") });
const migrationsDir = path.join(__dirname, "..", "database/migrations");
const entitiesDir = path.join(__dirname, "..", "database/entities");

const typeormConfig: DataSourceOptions = {
  type: "postgres",
  host: process.env.POSTGRES_HOST,
  port: Number(process.env.POSTGRES_PORT) || 40001,
  username: process.env.POSTGRES_USER,
  password: process.env.POSTGRES_PASSWORD,
  database: process.env.POSTGRES_DB,
  entities: [path.join(entitiesDir, "/**/*{.ts,.js}")],
  synchronize: false,
  migrationsRun: process.env.POSTGRES_MIGRATIONS_RUN === "true",
  migrations: [path.join(migrationsDir, "/**/*{.ts,.js}")],
  logging: process.env.LOGGING === "true",
};

export default typeormConfig;

export const AppDataSource = new DataSource(typeormConfig);
