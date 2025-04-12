package mockbuyticketstxrep

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockBuyTicketsTxRep struct {
	mock.Mock
}

func (m *MockBuyTicketsTxRep) Get(key uuid.UUID) (*models.BuyTicketTx, bool) {
	args := m.Called(key)
	return args.Get(0).(*models.BuyTicketTx), args.Bool(1)
}

func (m *MockBuyTicketsTxRep) Put(key uuid.UUID, value *models.BuyTicketTx) {
	m.Called(key, value)
}

func (m *MockBuyTicketsTxRep) Delete(key uuid.UUID) {
	m.Called(key)
}
