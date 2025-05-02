package ticketpurchasesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
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

func NewTicketPurchasesRep() (TicketPurchasesRep, error) {
	return &MockTicketPurchasesRep{}, nil
}
