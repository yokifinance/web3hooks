import { createParamDecorator, ExecutionContext } from "@nestjs/common";
import { AuthClient } from "./dto/auth.client";

export const Client = createParamDecorator(
  (data: unknown, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    return request.client as AuthClient;
  }
);

export type Client<Prop extends keyof AuthClient | undefined = undefined> =
  Prop extends keyof AuthClient ? AuthClient[Prop] : AuthClient;
