import { DynamicModule, Module } from "@nestjs/common";
import { TypeOrmModule, TypeOrmModuleOptions } from "@nestjs/typeorm";
import { ConfigService } from "@nestjs/config";
import * as path from "path";
import { EntitySchema } from "typeorm";
import { BaseConfig } from "../config";
import { ClientEntity } from "./entities/client.entity";
import { GueJobEntity } from "./entities/gueJob.entity";
import { SupportedChainEntity } from "./entities/supportedChain.entity";
import { EventListenerEntity } from "./entities/eventListener.entity";

const migrationsDir = path.join(__dirname, "migrations");
const entitiesDir = path.join(__dirname, "entities");
const migrations = [path.join(migrationsDir, "/**/*{.ts,.js}")];
//const entities = [path.join(entitiesDir, "/**/*{.ts,.js}")];

const entities = [
  ClientEntity,
  GueJobEntity,
  SupportedChainEntity,
  EventListenerEntity,
];

function typeOrmModulesFactory(
  appendOptions: Pick<
    Partial<TypeOrmModuleOptions>,
    | "entities"
    | "migrations"
    | "migrationsRun"
    | "migrationsTransactionMode"
    | "logger"
    | "migrationsTableName"
    | "metadataTableName"
  > = {}
) {
  return [
    TypeOrmModule.forFeature(entities),
    TypeOrmModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (
        configService: ConfigService<BaseConfig>
      ): TypeOrmModuleOptions => {
        return {
          ...configService.get("postgress"),
          ...appendOptions,
        };
      },
    }),
  ];
}

export type DataBaseModuleConfiguration = {
  migrationsTransactionMode: "all" | "none" | "each" | undefined;
  migrationsTableName: string;
  metadataTableName: string;
};

@Module({
  exports: [TypeOrmModule],
})
export class DatabaseModule {
  static forRoot(
    options: DataBaseModuleConfiguration | undefined = undefined
  ): DynamicModule {
    return {
      global: true,
      module: DatabaseModule,
      imports: [
        ...typeOrmModulesFactory({
          entities,
          migrations,
          migrationsTransactionMode:
            options?.migrationsTransactionMode || "each",
          migrationsTableName: options?.migrationsTableName || "migrations",
          metadataTableName:
            options?.metadataTableName || "new_typeorm_metadata",
        }),
      ],
    };
  }

  static forFeature(entities: EntitySchema[] = []): DynamicModule {
    return {
      global: true,
      module: DatabaseModule,
      imports: [
        ...typeOrmModulesFactory({
          entities,
          logger: "advanced-console",
        }),
      ],
    };
  }
}
