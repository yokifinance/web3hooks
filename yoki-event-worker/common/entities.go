package common

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type EventListener struct {
	bun.BaseModel `bun:"table:event_listeners"`

	Id         uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Chain      int
	Address    string
	ClientId   uuid.UUID
	WebhookUrl string
}

type LastProcessedBlock struct {
	bun.BaseModel `bun:"table:last_processed_blocks"`

	Chain                    int
	LastProcessedBlockNumber bunbig.Int
}
