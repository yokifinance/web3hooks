import { PostgresConnectionOptions } from "typeorm/driver/postgres/PostgresConnectionOptions";

export type DBConfig = PostgresConnectionOptions;

const getDefaultConfig = () =>
  ({
    type: "postgres",
    host: process.env.POSTGRES_HOST,
    port: Number(process.env.POSTGRES_PORT),
    username: process.env.POSTGRES_USER,
    database: process.env.POSTGRES_DB,
    password: process.env.POSTGRES_PASSWORD,
    logging: process.env.POSTGRES_LOG === "true",
    migrationsRun: process.env.POSTGRES_MIGRATION_RUN === "true",
  }) as DBConfig;

export default getDefaultConfig;
