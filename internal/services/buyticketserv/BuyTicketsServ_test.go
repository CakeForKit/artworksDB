package buyticketserv_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/buyticketserv"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuyTicketsServ_BuyTicket(t *testing.T) {
	t.Parallel()
	appCnfgCreator := testobj.NewAppConfigMother()
	eventCreator := testobj.NewEventMother()
	userCreator := testobj.NewUserMother()
	tptxCreator := testobj.NewTicketPurchaseTxMother()

	t.Run("success with authenticated user", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := appCnfgCreator.Default()
		event := eventCreator.EventCntTicketsP(uuid.New(), 100)
		user := userCreator.DefaultUserP(uuid.New())
		txCnt := 4
		purchasesCnt := 2
		expectedTX := tptxCreator.TicketPurchaseTxByUserP(uuid.New(), user.GetUsername(), user.GetEmail(), event.GetID(), user.GetID(), purchasesCnt)

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, event.GetID()).Return(event, nil)
		mockTXRep.On("GetCntTxByEventID", ctx, event.GetID()).Return(txCnt, nil)
		mockTPurchasesRep.On("GetCntTPurchasesForEvent", ctx, event.GetID()).Return(purchasesCnt, nil)
		mockAuthZ.On("UserIDFromContext", ctx).Return(user.GetID(), nil)
		mockUsrRep.On("GetByID", ctx, user.GetID()).Return(user, nil)
		mockTXRep.On("Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx")).Return(nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		tx, err := buyTicketsServ.BuyTicket(ctx, event.GetID(), purchasesCnt, "cn", "ce")

		require.Nil(t, err)
		require.NotNil(t, tx)
		require.True(t, expectedTX.Equals(tx))
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockUsrRep.AssertCalled(t, "GetByID", ctx, user.GetID())
		mockTXRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx"))
	})
	t.Run("success with unauthenticated user", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := appCnfgCreator.Default()
		event := eventCreator.EventCntTicketsP(uuid.New(), 100)
		txCnt := 4
		purchasesCnt := 2
		customerName := "Test Customer"
		customerEmail := "test@example.com"

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, event.GetID()).Return(event, nil)
		mockTXRep.On("GetCntTxByEventID", ctx, event.GetID()).Return(txCnt, nil)
		mockTPurchasesRep.On("GetCntTPurchasesForEvent", ctx, event.GetID()).Return(purchasesCnt, nil)
		mockAuthZ.On("UserIDFromContext", ctx).Return(uuid.Nil, auth.ErrNotAuthZ)
		mockTXRep.On("Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx")).Return(nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		tx, err := buyTicketsServ.BuyTicket(ctx, event.GetID(), purchasesCnt, customerName, customerEmail)

		require.Nil(t, err)
		require.NotNil(t, tx)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockUsrRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		mockTXRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx"))
	})

	t.Run("error no free tickets", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := appCnfgCreator.Default()
		event := eventCreator.EventCntTicketsP(uuid.New(), 5)
		txCnt := 3
		purchasesCnt := 3 // No free tickets (5 - 3 - 3 = -1)

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, event.GetID()).Return(event, nil)
		mockTXRep.On("GetCntTxByEventID", ctx, event.GetID()).Return(txCnt, nil)
		mockTPurchasesRep.On("GetCntTPurchasesForEvent", ctx, event.GetID()).Return(purchasesCnt, nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		_, err = buyTicketsServ.BuyTicket(ctx, event.GetID(), 1, "cn", "ce")

		require.Error(t, err)
		require.ErrorIs(t, err, buyticketserv.ErrNoFreeTicket)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertNotCalled(t, "UserIDFromContext", mock.Anything)
	})

	t.Run("error no user data for unauthenticated user", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := appCnfgCreator.Default()
		event := eventCreator.EventCntTicketsP(uuid.New(), 100)
		txCnt := 4
		purchasesCnt := 2

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, event.GetID()).Return(event, nil)
		mockTXRep.On("GetCntTxByEventID", ctx, event.GetID()).Return(txCnt, nil)
		mockTPurchasesRep.On("GetCntTPurchasesForEvent", ctx, event.GetID()).Return(purchasesCnt, nil)
		mockAuthZ.On("UserIDFromContext", ctx).Return(uuid.Nil, auth.ErrNotAuthZ)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		_, err = buyTicketsServ.BuyTicket(ctx, event.GetID(), 1, "", "")

		require.Error(t, err)
		require.ErrorIs(t, err, buyticketserv.ErrNoUserData)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
	})
}

func TestBuyTicketsServ_ConfirmBuyTicket(t *testing.T) {
	t.Parallel()
	tptxCreator := testobj.NewTicketPurchaseTxMother()

	t.Run("success confirm buy ticket", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		txID := uuid.New()
		tx := tptxCreator.TicketPurchaseTxP(txID, uuid.New(), uuid.New(), 1)

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockTXRep.On("GetByID", ctx, txID).Return(tx, nil)
		mockTPurchasesRep.On("Add", ctx, tx.GetTicketPurchase()).Return(nil)
		mockTXRep.On("Delete", ctx, txID).Return(nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		err = buyTicketsServ.ConfirmBuyTicket(ctx, txID)

		require.Nil(t, err)
		mockTXRep.AssertCalled(t, "GetByID", ctx, txID)
		mockTPurchasesRep.AssertCalled(t, "Add", ctx, tx.GetTicketPurchase())
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})

	t.Run("error transaction not found", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		txID := uuid.New()
		expectedErr := errors.New("transaction not found")

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockTXRep.On("GetByID", ctx, txID).Return(nil, expectedErr)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		err = buyTicketsServ.ConfirmBuyTicket(ctx, txID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "ConfirmBuyTicket")
		mockTXRep.AssertCalled(t, "GetByID", ctx, txID)
		mockTPurchasesRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

}

func TestBuyTicketsServ_CancelBuyTicket(t *testing.T) {
	t.Parallel()
	t.Run("success cancel buy ticket", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		txID := uuid.New()

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockTXRep.On("Delete", ctx, txID).Return(nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		err = buyTicketsServ.CancelBuyTicket(ctx, txID)

		require.Nil(t, err)
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})

	t.Run("error in cancel buy ticket", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		txID := uuid.New()
		expectedErr := errors.New("delete error")

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockTXRep.On("Delete", ctx, txID).Return(expectedErr)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		err = buyTicketsServ.CancelBuyTicket(ctx, txID)

		require.ErrorIs(t, err, expectedErr)
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})
}

func TestBuyTicketsServ_GetAllTicketPurchasesOfUser(t *testing.T) {
	t.Parallel()
	ticketPurchaseCreator := testobj.NewTicketPurchaseMother()

	t.Run("success get all ticket purchases of user", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		userID := uuid.New()
		tPurchases := []*models.TicketPurchase{
			ticketPurchaseCreator.TicketPurchaseP(uuid.New(), uuid.New(), userID),
			ticketPurchaseCreator.TicketPurchaseP(uuid.New(), uuid.New(), userID),
		}

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockAuthZ.On("UserIDFromContext", ctx).Return(userID, nil)
		mockTPurchasesRep.On("GetTPurchasesOfUserID", ctx, userID).Return(tPurchases, nil)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		resPurchases, err := buyTicketsServ.GetAllTicketPurchasesOfUser(ctx)

		require.Nil(t, err)
		require.True(t, len(tPurchases) == len(resPurchases))
		for i := range len(resPurchases) {
			require.True(t, tPurchases[i].Equals(resPurchases[i]))
		}
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockTPurchasesRep.AssertCalled(t, "GetTPurchasesOfUserID", ctx, userID)
	})

	t.Run("error user not authenticated", func(t *testing.T) {
		ctx := context.Background()
		appCnfg := testobj.NewAppConfigMother().Default()
		expectedErr := auth.ErrNotAuthZ

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockAuthZ.On("UserIDFromContext", ctx).Return(uuid.Nil, expectedErr)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		_, err = buyTicketsServ.GetAllTicketPurchasesOfUser(ctx)

		require.Error(t, err)
		require.ErrorIs(t, err, buyticketserv.ErrBuyTicketsServ)
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockTPurchasesRep.AssertNotCalled(t, "GetTPurchasesOfUserID", mock.Anything, mock.Anything)
	})
}

func TestBuyTicketsServ_GetBuyTicketTransactionDuration(t *testing.T) {
	t.Parallel()
	t.Run("success get transaction duration", func(t *testing.T) {
		appCnfg := testobj.NewAppConfigMother().Default()
		expectedDuration := 10 * time.Minute
		appCnfg.BuyTicketTransactionDuration = expectedDuration

		mockTXRep := new(buyticketstxrep.MockBuyTicketsTxRep)
		mockTPurchasesRep := new(ticketpurchasesrep.MockTicketPurchasesRep)
		mockAuthZ := new(auth.MockAuthZ)
		mockUsrRep := new(userrep.MockUserRep)
		mockEventRep := new(eventrep.MockEventRep)

		buyTicketsServ, err := buyticketserv.NewBuyTicketsServ(
			mockTXRep,
			mockTPurchasesRep,
			appCnfg,
			mockAuthZ,
			mockUsrRep,
			mockEventRep,
		)
		require.Nil(t, err)

		// ACT
		duration := buyTicketsServ.GetBuyTicketTransactionDuration()

		require.Equal(t, expectedDuration, duration)
	})
}
