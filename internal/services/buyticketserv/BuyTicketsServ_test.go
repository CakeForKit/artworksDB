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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type BuyTicketsServiceSuite struct {
	suite.Suite
}

func TestBuyTicketsService(t *testing.T) {
	suite.RunSuite(t, new(BuyTicketsServiceSuite))
}

func (s *BuyTicketsServiceSuite) TestBuyTicketsServ_BuyTicket(t provider.T) {
	t.Parallel()
	appCnfgCreator := testobj.NewAppConfigMother()
	eventCreator := testobj.NewEventMother()
	userCreator := testobj.NewUserMother()
	tptxCreator := testobj.NewTicketPurchaseTxMother()

	t.WithNewStep("success with authenticated user", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		tx, err := buyTicketsServ.BuyTicket(ctx, event.GetID(), purchasesCnt, "cn", "ce")

		sCtx.Require().NoError(err)
		sCtx.Assert().NotNil(tx)
		sCtx.Require().True(expectedTX.Equals(tx))
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockUsrRep.AssertCalled(t, "GetByID", ctx, user.GetID())
		mockTXRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx"))
	})
	t.WithNewStep("success with unauthenticated user", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		tx, err := buyTicketsServ.BuyTicket(ctx, event.GetID(), purchasesCnt, customerName, customerEmail)

		sCtx.Require().NoError(err)
		sCtx.Assert().NotNil(tx)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockUsrRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		mockTXRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("models.TicketPurchaseTx"))
	})

	t.WithNewStep("error no free tickets", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		_, err = buyTicketsServ.BuyTicket(ctx, event.GetID(), 1, "cn", "ce")

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, buyticketserv.ErrNoFreeTicket)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertNotCalled(t, "UserIDFromContext", mock.Anything)
	})

	t.WithNewStep("error no user data for unauthenticated user", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		_, err = buyTicketsServ.BuyTicket(ctx, event.GetID(), 1, "", "")

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, buyticketserv.ErrNoUserData)
		mockEventRep.AssertCalled(t, "GetByID", ctx, event.GetID())
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
	})
}

func (s *BuyTicketsServiceSuite) TestBuyTicketsServ_ConfirmBuyTicket(t provider.T) {
	t.Parallel()
	tptxCreator := testobj.NewTicketPurchaseTxMother()

	t.WithNewStep("success confirm buy ticket", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		err = buyTicketsServ.ConfirmBuyTicket(ctx, txID)

		sCtx.Require().NoError(err)
		mockTXRep.AssertCalled(t, "GetByID", ctx, txID)
		mockTPurchasesRep.AssertCalled(t, "Add", ctx, tx.GetTicketPurchase())
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})

	t.WithNewStep("error transaction not found", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		err = buyTicketsServ.ConfirmBuyTicket(ctx, txID)

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockTXRep.AssertCalled(t, "GetByID", ctx, txID)
		mockTPurchasesRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

}

func (s *BuyTicketsServiceSuite) TestBuyTicketsServ_CancelBuyTicket(t provider.T) {
	t.Parallel()
	t.WithNewStep("success cancel buy ticket", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		err = buyTicketsServ.CancelBuyTicket(ctx, txID)

		sCtx.Require().NoError(err)
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})

	t.WithNewStep("error in cancel buy ticket", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		err = buyTicketsServ.CancelBuyTicket(ctx, txID)

		sCtx.Require().ErrorIs(err, expectedErr)
		mockTXRep.AssertCalled(t, "Delete", ctx, txID)
	})
}

func (s *BuyTicketsServiceSuite) TestBuyTicketsServ_GetAllTicketPurchasesOfUser(t provider.T) {
	t.Parallel()
	ticketPurchaseCreator := testobj.NewTicketPurchaseMother()

	t.WithNewStep("success get all ticket purchases of user", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		resPurchases, err := buyTicketsServ.GetAllTicketPurchasesOfUser(ctx)

		sCtx.Require().NoError(err)
		sCtx.Require().True(len(tPurchases) == len(resPurchases))
		for i := range len(resPurchases) {
			sCtx.Require().True(tPurchases[i].Equals(resPurchases[i]))
		}
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockTPurchasesRep.AssertCalled(t, "GetTPurchasesOfUserID", ctx, userID)
	})

	t.WithNewStep("error user not authenticated", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		_, err = buyTicketsServ.GetAllTicketPurchasesOfUser(ctx)

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, buyticketserv.ErrBuyTicketsServ)
		mockAuthZ.AssertCalled(t, "UserIDFromContext", ctx)
		mockTPurchasesRep.AssertNotCalled(t, "GetTPurchasesOfUserID", mock.Anything, mock.Anything)
	})
}

func (s *BuyTicketsServiceSuite) TestBuyTicketsServ_GetBuyTicketTransactionDuration(t provider.T) {
	t.Parallel()
	t.WithNewStep("success get transaction duration", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)

		// ACT
		duration := buyTicketsServ.GetBuyTicketTransactionDuration()

		sCtx.Assert().Equal(expectedDuration, duration)
	})
}
