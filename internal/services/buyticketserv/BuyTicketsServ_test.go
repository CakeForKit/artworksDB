package buyticketserv_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/buyticketserv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// testHelper содержит общие методы для тестов
type testHelper struct {
	ctx context.Context
	// serv       buyticketserv.BuyTicketsServ
	// txRep      *buyticketstxrep.MockBuyTicketsTxRep
	// tpRep      *ticketpurchasesrep.MockTicketPurchasesRep
	config     cnfg.AppConfig
	employeeID uuid.UUID
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	config := cnfg.AppConfig{
		BuyTicketTransactionDuration: 15 * time.Minute,
	}

	// txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
	// tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)

	// serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, config)
	// require.NoError(t, err)

	return &testHelper{
		ctx: ctx,
		// serv:       serv,
		// txRep:      txRep,
		// tpRep:      tpRep,
		config:     config,
		employeeID: uuid.New(),
	}
}

func (th *testHelper) createTestEvent(num int) *models.Event {
	event, err := models.NewEvent(
		uuid.New(),
		fmt.Sprintf("Event %d", num),
		time.Now().Add(time.Hour),
		time.Now().Add(2*time.Hour),
		fmt.Sprintf("Address %d", num),
		true,
		th.employeeID,
		100+num,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestEvent failed: %v", err))
	}
	return &event
}

func (th *testHelper) createTestUser(num int) *models.User {
	user, err := models.NewUser(
		uuid.New(),
		fmt.Sprintf("user%d", num),
		fmt.Sprintf("login%d", num),
		fmt.Sprintf("hash %d", num),
		time.Now(),
		fmt.Sprintf("user%d@example.com", num),
		true,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestUser failed: %v", err))
	}
	return &user
}

func (th *testHelper) createTestTicketPurchaseTx(eventID uuid.UUID, userID uuid.UUID, cnt int) *models.TicketPurchaseTx {
	tx, err := models.NewBuyTicketTx(
		uuid.New(),
		"Customer",
		"customer@example.com",
		time.Now(),
		eventID,
		userID,
		cnt,
		time.Now().Add(th.config.BuyTicketTransactionDuration),
	)
	if err != nil {
		panic(fmt.Sprintf("createTestTicketPurchaseTx failed: %v", err))
	}
	return &tx
}

func TestBuyTicketsServ_BuyTicket(t *testing.T) {
	th := setupTestHelper(t)

	event := th.createTestEvent(1)
	customerName := "Test Customer"
	customerEmail := "test@example.com"
	cntTickets := 2

	t.Run("Should successfully buy tickets", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		// Настраиваем моки
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("GetCntTPurchasesForEvent", th.ctx, event.GetID()).Return(0, nil)
		txRep.On("Add", th.ctx, mock.AnythingOfType("models.TicketPurchaseTx")).Return(nil)

		tx, err := serv.BuyTicket(th.ctx, *event, cntTickets, customerName, customerEmail)
		require.NoError(t, err)

		assert.Equal(t, cntTickets, tx.GetCntTickets())
		assert.Equal(t, customerName, tx.GetTicketPurchase().GetCustomerName())
		assert.Equal(t, customerEmail, tx.GetTicketPurchase().GetCustomerEmail())
		assert.Equal(t, event.GetID(), tx.GetTicketPurchase().GetEventID())
		assert.True(t, tx.GetExpiredAt().After(time.Now()))

		txRep.AssertExpectations(t)
	})

	t.Run("Should return error when no free tickets", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		// Настраиваем моки - все билеты заняты
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("GetCntTPurchasesForEvent", th.ctx, event.GetID()).Return(event.GetTicketCount(), nil)

		_, err = serv.BuyTicket(th.ctx, *event, cntTickets, customerName, customerEmail)
		assert.Error(t, err)                                                   // Проверяем что ошибка есть
		assert.Contains(t, err.Error(), buyticketserv.ErrNoFreeTicket.Error()) // Проверяем текст ошибки

		txRep.AssertExpectations(t)
		tpRep.AssertExpectations(t)
	})

	t.Run("Should return error when cntTickets <= 0", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)

		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("GetCntTPurchasesForEvent", th.ctx, event.GetID()).Return(event.GetTicketCount(), nil)

		_, err = serv.BuyTicket(th.ctx, *event, 0, customerName, customerEmail)
		assert.Error(t, err)
	})
}

func TestBuyTicketsServ_BuyTicketByUser(t *testing.T) {
	th := setupTestHelper(t)

	event := th.createTestEvent(1)
	user := th.createTestUser(1)
	cntTickets := 2

	t.Run("Should successfully buy tickets by user", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		// Настраиваем моки
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("GetCntTPurchasesForEvent", th.ctx, event.GetID()).Return(0, nil)
		txRep.On("Add", th.ctx, mock.AnythingOfType("models.TicketPurchaseTx")).Return(nil)

		tx, err := serv.BuyTicketByUser(th.ctx, *event, cntTickets, *user)
		require.NoError(t, err)

		assert.Equal(t, cntTickets, tx.GetCntTickets())
		assert.Equal(t, user.GetUsername(), tx.GetTicketPurchase().GetCustomerName())
		assert.Equal(t, user.GetEmail(), tx.GetTicketPurchase().GetCustomerEmail())
		assert.Equal(t, user.GetID(), tx.GetTicketPurchase().GetUserID())
		assert.Equal(t, event.GetID(), tx.GetTicketPurchase().GetEventID())
		assert.True(t, tx.GetExpiredAt().After(time.Now()))

		txRep.AssertExpectations(t)
	})

	t.Run("Should return error when no free tickets", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		// Настраиваем моки - все билеты заняты
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("GetCntTPurchasesForEvent", th.ctx, event.GetID()).Return(event.GetTicketCount(), nil)

		_, err = serv.BuyTicketByUser(th.ctx, *event, cntTickets, *user)
		assert.Error(t, err)                                                   // Проверяем что ошибка есть
		assert.Contains(t, err.Error(), buyticketserv.ErrNoFreeTicket.Error()) // Проверяем текст ошибки

		txRep.AssertExpectations(t)
		tpRep.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_ConfirmBuyTicket(t *testing.T) {
	th := setupTestHelper(t)

	event := th.createTestEvent(1)
	user := th.createTestUser(1)
	tx := th.createTestTicketPurchaseTx(event.GetID(), user.GetID(), 2)

	t.Run("Should successfully confirm ticket purchase", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		// Настраиваем моки
		txRep.On("GetByID", th.ctx, tx.GetID()).Return(tx, nil)
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(1, nil).Once()
		txRep.On("Delete", th.ctx, tx.GetID()).Return(nil)
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil).Once()
		tpRep.On("Add", th.ctx, tx.GetTicketPurchase()).Return(nil)

		err = serv.ConfirmBuyTicket(th.ctx, tx.GetID())
		require.NoError(t, err)

		txRep.AssertExpectations(t)
		tpRep.AssertExpectations(t)
	})

	t.Run("Should return error when transaction not found", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		nonExistentID := uuid.New()
		txRep.On("GetByID", th.ctx, nonExistentID).Return(nil, errors.New("not found"))

		err = serv.ConfirmBuyTicket(th.ctx, nonExistentID)
		assert.Error(t, err)

		txRep.AssertExpectations(t)
	})

	t.Run("Should return error when failed to delete transaction", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		txRep.On("GetByID", th.ctx, tx.GetID()).Return(tx, nil)
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(1, nil)
		txRep.On("Delete", th.ctx, tx.GetID()).Return(errors.New("delete error"))

		err = serv.ConfirmBuyTicket(th.ctx, tx.GetID())
		assert.Error(t, err)

		txRep.AssertExpectations(t)
	})

	t.Run("Should return error when failed to add ticket purchase", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		txRep.On("GetByID", th.ctx, tx.GetID()).Return(tx, nil)
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(1, nil)
		txRep.On("Delete", th.ctx, tx.GetID()).Return(nil)
		txRep.On("GetCntTxByEventID", th.ctx, event.GetID()).Return(0, nil)
		tpRep.On("Add", th.ctx, tx.GetTicketPurchase()).Return(errors.New("add error"))

		err = serv.ConfirmBuyTicket(th.ctx, tx.GetID())
		assert.Error(t, err)

		txRep.AssertExpectations(t)
		tpRep.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_CancelBuyTicket(t *testing.T) {
	th := setupTestHelper(t)

	txID := uuid.New()

	t.Run("Should successfully cancel ticket purchase", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		txRep.On("Delete", th.ctx, txID).Return(nil)

		err = serv.CancelBuyTicket(th.ctx, txID)
		require.NoError(t, err)

		txRep.AssertExpectations(t)
	})

	t.Run("Should return error when failed to cancel", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		txRep.On("Delete", th.ctx, txID).Return(errors.New("delete error"))

		err = serv.CancelBuyTicket(th.ctx, txID)
		assert.Error(t, err)

		txRep.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_GetAllTicketPurchasesOfUser(t *testing.T) {
	th := setupTestHelper(t)

	user := th.createTestUser(1)
	event1 := th.createTestEvent(1)
	event2 := th.createTestEvent(2)

	tx1 := th.createTestTicketPurchaseTx(event1.GetID(), user.GetID(), 1)
	tx2 := th.createTestTicketPurchaseTx(event2.GetID(), user.GetID(), 2)

	t.Run("Should return all ticket purchases of user", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		expectedPurchases := []*models.TicketPurchase{
			tx1.GetTicketPurchase(),
			tx2.GetTicketPurchase(),
		}

		tpRep.On("GetTPurchasesOfUserID", th.ctx, user.GetID()).Return(expectedPurchases, nil)

		purchases, err := serv.GetAllTicketPurchasesOfUser(th.ctx, user.GetID())
		require.NoError(t, err)

		assert.Len(t, purchases, 2)
		assert.Equal(t, tx1.GetTicketPurchase().GetID(), purchases[0].GetID())
		assert.Equal(t, tx2.GetTicketPurchase().GetID(), purchases[1].GetID())

		tpRep.AssertExpectations(t)
	})

	t.Run("Should return error when failed to get purchases", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		tpRep.On("GetTPurchasesOfUserID", th.ctx, user.GetID()).Return(make([]*models.TicketPurchase, 0), errors.New("get error"))

		_, err = serv.GetAllTicketPurchasesOfUser(th.ctx, user.GetID())
		assert.Error(t, err)

		tpRep.AssertExpectations(t)
	})

	t.Run("Should return empty list when no purchases", func(t *testing.T) {
		txRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		tpRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		serv, err := buyticketserv.NewBuyTicketsServ(txRep, tpRep, th.config)
		require.NoError(t, err)
		tpRep.On("GetTPurchasesOfUserID", th.ctx, user.GetID()).Return([]*models.TicketPurchase{}, nil)

		purchases, err := serv.GetAllTicketPurchasesOfUser(th.ctx, user.GetID())
		require.NoError(t, err)

		assert.Empty(t, purchases)

		tpRep.AssertExpectations(t)
	})
}
