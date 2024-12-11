import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { BaseConfig } from "common/config";
import { AuthClient } from "./dto/auth.client";
import { InjectRepository } from "@nestjs/typeorm";
import { ClientEntity } from "common/database/entities/client.entity";
import { Repository } from "typeorm";
import { Reflector } from "@nestjs/core";
import { IS_PUBLIC_KEY } from "./public.decorator";

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(
    private reflector: Reflector,
    @InjectRepository(ClientEntity)
    private clientRepository: Repository<ClientEntity>
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) {
      return true;
    }

    const request = context.switchToHttp().getRequest();
    const secretKey = request.headers["secret-key"];
    if (!secretKey) {
      throw new UnauthorizedException();
    }

    try {
      const client = await this.clientRepository.findOneBy({ secretKey });
      if (!client) throw new UnauthorizedException();

      request["client"] = { ...client } as AuthClient;
    } catch {
      throw new UnauthorizedException();
    }
    return true;
  }
}
