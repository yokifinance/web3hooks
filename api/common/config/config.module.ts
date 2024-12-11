import { ConfigModule } from "@nestjs/config";
import getDefaultApplicationConfig, {
  ApplicationConfig,
} from "./application.config";
import getDefaultDataBaseConfig, { DBConfig } from "./db.config";

export type BaseConfig = {
  application: ApplicationConfig;
  postgress: DBConfig;
};

export type CustomConfig = {
  [section: string]: Object;
};

const loadCustomConfig = <T extends CustomConfig>(
  customConfig: T
): BaseConfig & T => {
  const application = getDefaultApplicationConfig();
  const postgress = getDefaultDataBaseConfig();
  // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
  /*
    For simplicity - just add new app configs here
    Pay attention that they have to be functions since "process.env" is not available outside of config.module

    if desired - it can still be accessed here, ex.
    const rps = { executorAddress: process.env.ADDRESS }
  */
  // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
  return {
    application,
    postgress,
    ...customConfig,
  };
};

export const getGlobalConfigModule = <T extends CustomConfig>(
  customConfig: T
) => {
  const loadConfig = () => {
    return loadCustomConfig(customConfig);
  };
  return ConfigModule.forRoot({
    isGlobal: true,
    load: [loadConfig],
    envFilePath: [".env", "../.env", "/run/secrets/yoki_web3tasks_env"],
  });
};
