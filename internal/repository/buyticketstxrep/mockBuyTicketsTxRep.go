package buyticketstxrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockBuyTicketsTxRep struct {
	mock.Mock
}

func (m *MockBuyTicketsTxRep) GetByID(ctx context.Context, txID uuid.UUID) (*models.TicketPurchaseTx, error) {
	args := m.Called(ctx, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TicketPurchaseTx), args.Error(1)
}

func (m *MockBuyTicketsTxRep) GetCntTxByEventID(ctx context.Context, eventID uuid.UUID) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *MockBuyTicketsTxRep) Add(ctx context.Context, tpTx models.TicketPurchaseTx) error {
	args := m.Called(ctx, tpTx)
	return args.Error(0)
}

func (m *MockBuyTicketsTxRep) Delete(ctx context.Context, txID uuid.UUID) error {
	args := m.Called(ctx, txID)
	return args.Error(0)
}

func (m *MockBuyTicketsTxRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBuyTicketsTxRep) Close() {
	m.Called()
}
