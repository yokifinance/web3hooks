# Development

## Starting

`npm i`
`npm run start:dev`
`vscode -> debug -> yoki-api`
`docker-compose up`

To start only db:
`docker-compose up db`

# Migrations

Generate
`npm run typeorm migration:generate -- -d ./common/database/typeorm.config.ts ./common/database/migrations/Init`

Run
`npm run typeorm migration:run -- -d ./common/database/typeorm.config.ts`

Revert
`npm run typeorm migration:revert -- -d ./common/database/typeorm.config.ts`
