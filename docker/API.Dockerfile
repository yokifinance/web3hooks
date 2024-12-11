FROM node:21-alpine as build

WORKDIR /app

# RUN apk update && \
#     apk add --no-cache git=2.40.1-r0 && \
#     apk add --no-cache python3 && \
#     apk add --no-cache build-base

COPY ./api .

RUN npm ci
RUN npm run build

FROM node:21-alpine
WORKDIR /app

# TODO: build everything in /dist without node_modules
COPY --from=build /app/package.json /app/package.json
COPY --from=build /app/package-lock.json /app/package-lock.json
COPY --from=build /app/dist /app/dist
COPY --from=build /app/node_modules /app/node_modules

ENV APPLICATION_PORT=8001
EXPOSE 8001
CMD ["npm", "run", "start:prod"]
