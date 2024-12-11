// main params related to application (ex. port)
export enum ApplicationMode {
  dev = 'dev',
  prod = 'prod',
}

export type ApplicationConfig = {
  port: number;
  mode: ApplicationMode;
};

const getDefaultConfig = () =>
  ({
    port: Number(process.env.APPLICATION_PORT) || 8081,
    mode:
      process.env.APPLICATION_MODE === ApplicationMode.prod
        ? ApplicationMode.prod
        : ApplicationMode.dev,
  }) as ApplicationConfig;

export default getDefaultConfig;
