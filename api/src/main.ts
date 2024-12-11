import { Logger } from '@nestjs/common';
import bootstrap from './bootstrap';

const logger = new Logger('NestApplication', { timestamp: true });

bootstrap(logger).catch((error: unknown) => {
  logger.error(`API bootstrapping application failed! ${error}`);
});
