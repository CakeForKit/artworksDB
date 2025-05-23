package buyticketserv_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/buyticketserv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testData struct {
	ctx     context.Context
	config  cnfg.AppConfig
	userID  uuid.UUID
	eventID uuid.UUID
}

func setupTestData() *testData {
	return &testData{
		ctx:     context.Background(),
		config:  cnfg.AppConfig{BuyTicketTransactionDuration: 15 * time.Minute},
		userID:  uuid.New(),
		eventID: uuid.New(),
	}
}

func createTestEvent(eventID uuid.UUID, ticketCount int) *models.Event {
	event, _ := models.NewEvent(
		eventID,
		"Test Event",
		time.Now(),
		time.Now().Add(24*time.Hour),
		"Test Address",
		true,
		uuid.New(),
		ticketCount,
		true,
		make(uuid.UUIDs, 0),
	)
	return &event
}

func createTestUser(userID uuid.UUID) *models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login",
		"hashed-password",
		time.Now(),
		"user@test.com",
		true,
	)
	return &user
}

func createTestTicketPurchaseTx(eventID, userID uuid.UUID, config cnfg.AppConfig, cnt int) *models.TicketPurchaseTx {
	tx, _ := models.NewBuyTicketTx(
		uuid.New(),
		"Customer",
		"customer@example.com",
		time.Now(),
		eventID,
		userID,
		cnt,
		time.Now().Add(config.BuyTicketTransactionDuration),
	)
	return &tx
}

func TestBuyTicketsServ_BuyTicket(t *testing.T) {
	td := setupTestData()
	event := createTestEvent(td.eventID, 10)
	user := createTestUser(td.userID)
	customerName := "Test Customer"
	customerEmail := "test@example.com"
	cntTickets := 2

	t.Run("success for authenticated user", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		userMock := new(userrep.MockUserRep)
		eventMock := new(eventrep.MockEventRep)
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		authMock.On("UserIDFromContext", td.ctx).Return(td.userID, nil)
		userMock.On("GetByID", td.ctx, td.userID).Return(user, nil)
		eventMock.On("GetByID", td.ctx, td.eventID).Return(event, nil)
		txMock.On("GetCntTxByEventID", td.ctx, td.eventID).Return(0, nil)
		ticketMock.On("GetCntTPurchasesForEvent", td.ctx, td.eventID).Return(0, nil)
		txMock.On("Add", td.ctx, mock.Anything).Return(nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			ticketMock,
			td.config,
			authMock,
			userMock,
			eventMock,
		)
		require.NoError(t, err)

		tx, err := service.BuyTicket(td.ctx, td.eventID, cntTickets, "", "")
		require.NoError(t, err)

		assert.Equal(t, cntTickets, tx.GetCntTickets())
		assert.Equal(t, user.GetUsername(), tx.GetTicketPurchase().GetCustomerName())
		assert.Equal(t, user.GetEmail(), tx.GetTicketPurchase().GetCustomerEmail())
		assert.Equal(t, td.eventID, tx.GetTicketPurchase().GetEventID())
		assert.True(t, tx.GetExpiredAt().After(time.Now()))

		authMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
		eventMock.AssertExpectations(t)
		txMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})

	t.Run("success for unauthenticated user", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		eventMock := new(eventrep.MockEventRep)
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		authMock.On("UserIDFromContext", td.ctx).Return(uuid.Nil, auth.ErrNotAuthZ)
		eventMock.On("GetByID", td.ctx, td.eventID).Return(event, nil)
		txMock.On("GetCntTxByEventID", td.ctx, td.eventID).Return(0, nil)
		ticketMock.On("GetCntTPurchasesForEvent", td.ctx, td.eventID).Return(0, nil)
		txMock.On("Add", td.ctx, mock.Anything).Return(nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			ticketMock,
			td.config,
			authMock,
			new(userrep.MockUserRep),
			eventMock,
		)
		require.NoError(t, err)

		tx, err := service.BuyTicket(td.ctx, td.eventID, cntTickets, customerName, customerEmail)
		require.NoError(t, err)

		assert.Equal(t, cntTickets, tx.GetCntTickets())
		assert.Equal(t, customerName, tx.GetTicketPurchase().GetCustomerName())
		assert.Equal(t, customerEmail, tx.GetTicketPurchase().GetCustomerEmail())
		assert.Equal(t, td.eventID, tx.GetTicketPurchase().GetEventID())

		authMock.AssertExpectations(t)
		eventMock.AssertExpectations(t)
		txMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})

	t.Run("error when no free tickets", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		userMock := new(userrep.MockUserRep)
		eventMock := new(eventrep.MockEventRep)
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		// authMock.On("UserIDFromContext", td.ctx).Return(td.userID, nil)
		// userMock.On("GetByID", td.ctx, td.userID).Return(user, nil)
		eventMock.On("GetByID", td.ctx, td.eventID).Return(event, nil)
		txMock.On("GetCntTxByEventID", td.ctx, td.eventID).Return(8, nil)
		ticketMock.On("GetCntTPurchasesForEvent", td.ctx, td.eventID).Return(2, nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			ticketMock,
			td.config,
			authMock,
			userMock,
			eventMock,
		)
		require.NoError(t, err)

		_, err = service.BuyTicket(td.ctx, td.eventID, cntTickets, "", "")
		assert.ErrorIs(t, err, buyticketserv.ErrNoFreeTicket)

		// authMock.AssertExpectations(t)
		// userMock.AssertExpectations(t)
		eventMock.AssertExpectations(t)
		txMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})

	t.Run("error when no user data for unauthenticated user", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		eventMock := new(eventrep.MockEventRep)
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		// Set up mock expectations
		authMock.On("UserIDFromContext", td.ctx).Return(uuid.Nil, auth.ErrNotAuthZ)
		eventMock.On("GetByID", td.ctx, td.eventID).Return(event, nil)
		txMock.On("GetCntTxByEventID", td.ctx, td.eventID).Return(0, nil)
		ticketMock.On("GetCntTPurchasesForEvent", td.ctx, td.eventID).Return(0, nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			ticketMock,
			td.config,
			authMock,
			new(userrep.MockUserRep),
			eventMock,
		)
		require.NoError(t, err)

		_, err = service.BuyTicket(td.ctx, td.eventID, cntTickets, "", "")
		assert.ErrorIs(t, err, buyticketserv.ErrNoUserData)

		// Verify expected calls were made
		authMock.AssertExpectations(t)
		eventMock.AssertExpectations(t)
		txMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_ConfirmBuyTicket(t *testing.T) {
	td := setupTestData()
	tx := createTestTicketPurchaseTx(td.eventID, td.userID, td.config, 2)

	t.Run("success", func(t *testing.T) {
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		txMock.On("GetByID", td.ctx, tx.GetID()).Return(tx, nil)
		ticketMock.On("Add", td.ctx, tx.GetTicketPurchase()).Return(nil)
		txMock.On("Delete", td.ctx, tx.GetID()).Return(nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			ticketMock,
			td.config,
			new(auth.MockAuthZ),
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		err = service.ConfirmBuyTicket(td.ctx, tx.GetID())
		assert.NoError(t, err)

		txMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})

	t.Run("error when tx not found", func(t *testing.T) {
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		txMock.On("GetByID", td.ctx, tx.GetID()).Return(nil, errors.New("not found"))

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			new(ticketpurchasesrep.MockTicketPurchasesRep),
			td.config,
			new(auth.MockAuthZ),
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		err = service.ConfirmBuyTicket(td.ctx, tx.GetID())
		assert.Error(t, err)

		txMock.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_CancelBuyTicket(t *testing.T) {
	td := setupTestData()
	txID := uuid.New()

	t.Run("success", func(t *testing.T) {
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		txMock.On("Delete", td.ctx, txID).Return(nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			new(ticketpurchasesrep.MockTicketPurchasesRep),
			td.config,
			new(auth.MockAuthZ),
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		err = service.CancelBuyTicket(td.ctx, txID)
		assert.NoError(t, err)

		txMock.AssertExpectations(t)
	})

	t.Run("error when delete fails", func(t *testing.T) {
		txMock := new(buyticketstxrep.MockBuyTicketsTxRep)
		txMock.On("Delete", td.ctx, txID).Return(errors.New("delete error"))

		service, err := buyticketserv.NewBuyTicketsServ(
			txMock,
			new(ticketpurchasesrep.MockTicketPurchasesRep),
			td.config,
			new(auth.MockAuthZ),
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		err = service.CancelBuyTicket(td.ctx, txID)
		assert.Error(t, err)

		txMock.AssertExpectations(t)
	})
}

func TestBuyTicketsServ_GetAllTicketPurchasesOfUser(t *testing.T) {
	td := setupTestData()
	tx := createTestTicketPurchaseTx(td.eventID, td.userID, td.config, 1)
	purchases := []*models.TicketPurchase{tx.GetTicketPurchase()}

	t.Run("success", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		ticketMock := new(ticketpurchasesrep.MockTicketPurchasesRep)

		authMock.On("UserIDFromContext", td.ctx).Return(td.userID, nil)
		ticketMock.On("GetTPurchasesOfUserID", td.ctx, td.userID).Return(purchases, nil)

		service, err := buyticketserv.NewBuyTicketsServ(
			new(buyticketstxrep.MockBuyTicketsTxRep),
			ticketMock,
			td.config,
			authMock,
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		result, err := service.GetAllTicketPurchasesOfUser(td.ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, tx.GetTicketPurchase().GetID(), result[0].GetID())

		authMock.AssertExpectations(t)
		ticketMock.AssertExpectations(t)
	})

	t.Run("error when not authenticated", func(t *testing.T) {
		authMock := new(auth.MockAuthZ)
		authMock.On("UserIDFromContext", td.ctx).Return(uuid.Nil, auth.ErrNotAuthZ)

		service, err := buyticketserv.NewBuyTicketsServ(
			new(buyticketstxrep.MockBuyTicketsTxRep),
			new(ticketpurchasesrep.MockTicketPurchasesRep),
			td.config,
			authMock,
			new(userrep.MockUserRep),
			new(eventrep.MockEventRep),
		)
		require.NoError(t, err)

		_, err = service.GetAllTicketPurchasesOfUser(td.ctx)
		assert.Error(t, err)

		authMock.AssertExpectations(t)
	})
}
