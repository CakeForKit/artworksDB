package buyticketstxrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type BuyTicketsTxRep interface {
	GetByID(ctx context.Context, txID uuid.UUID) (*models.TicketPurchaseTx, error)
	GetCntTxByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
	Add(ctx context.Context, tpTx models.TicketPurchaseTx) error
	Delete(ctx context.Context, txID uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrExpireTx   = errors.New("transaction already expired")
	ErrTxNotFound = errors.New("transaction not found")
)

func NewBuyTicketsTxRep(
	ctx context.Context,
	redisCreds *cnfg.RedisCredentials,
) (BuyTicketsTxRep, error) {
	return &MockBuyTicketsTxRep{}, nil
	// return NewRedisBuyTicketsTxRep(ctx, redisCreds)
}
