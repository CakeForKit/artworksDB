package ticketpurchasesrep

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrPgTicketPurchasesRep = errors.New("pgTicketPurchasesRep")
)

type TicketPurchasesRep interface {
	// GetByID(ctx context.Context, id uuid.UUID) (*models.TicketPurchase, error)
	GetTPurchasesOfUserID(ctx context.Context, userID uuid.UUID) ([]*models.TicketPurchase, error)
	GetCntTPurchasesForEvent(ctx context.Context, eventID uuid.UUID) (int, error)
	Add(ctx context.Context, tp *models.TicketPurchase) error
	// Delete(ctx context.Context, id uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

func NewTicketPurchasesRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (TicketPurchasesRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgTicketPurchasesRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHTicketPurchasesRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewTicketPurchasesRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
