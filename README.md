# Yoki-web3hooks

Ever needed to receive webhooks for blockchain events on EVM chains? We did! Existing commercial solutions can be
too expensive for an early-stage startup. Also, what if you want full control of what's happening?
We encountered this task in Yoki Finance and developed this open-source library
to help you with webhooks – welcome **Yoki-web3hooks**!

## Project folders

- `api` - API written with Nest.js to add webhook listeners.
- `yoki-event-worker` - worker that listens to blockchain events and puts webhooks into the queue (written in Golang)
- `yoki-webhook-executor` - worker that delivers webhooks to listeners (written in Golang)

Only Postgres DB is needed.

## Development

Rename `.env.example` -> `.env`
If you use tests, to `test.env`

- `cd /api`
- `npm i`
- `npm run start:dev` or `F5 in vscode`

Manually create "supported_chain" in db

## Testing

Clear go test cache
`go clean -testcache`

## Deploy

Deployment secrets:

- `/run/secrets/yoki_web3tasks_env` – Docker secrets if you want to use them in your pipeline.
