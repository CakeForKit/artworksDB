package ticketpurchasesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTicketPurchasesRep struct {
	mock.Mock
}

func (m *MockTicketPurchasesRep) GetTPurchasesOfUserID(ctx context.Context, userID uuid.UUID) ([]*models.TicketPurchase, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.TicketPurchase), args.Error(1)
}

func (m *MockTicketPurchasesRep) GetCntTPurchasesForEvent(ctx context.Context, eventID uuid.UUID) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *MockTicketPurchasesRep) Add(ctx context.Context, tp *models.TicketPurchase) error {
	args := m.Called(ctx, tp)
	return args.Error(0)
}

func (m *MockTicketPurchasesRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTicketPurchasesRep) Close() {
	m.Called()
}
