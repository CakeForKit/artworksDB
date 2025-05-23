package buyticketserv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/google/uuid"
)

var (
	ErrBuyTicketsServ = errors.New("buyTicketsServ")
	ErrNoFreeTicket   = errors.New("no free ticket")
	ErrNoUserData     = errors.New("no info about user (customerName, customerEmail)")
)

type BuyTicketsServ interface {
	BuyTicket(ctx context.Context, eventID uuid.UUID, cntTickets int, customerName string, customerEmail string) (*models.TicketPurchaseTx, error)
	// BuyTicketByUser(ctx context.Context, event models.Event, cntTickets int, user models.User) (*models.TicketPurchaseTx, error)
	ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error
	CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error
	GetAllTicketPurchasesOfUser(ctx context.Context) ([]*models.TicketPurchase, error)
	GetBuyTicketTransactionDuration() time.Duration
}

type buyTicketsServ struct {
	txRep         buyticketstxrep.BuyTicketsTxRep
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep
	config        cnfg.AppConfig
	authZ         auth.AuthZ
	userRep       userrep.UserRep
	eventRep      eventrep.EventRep
}

func NewBuyTicketsServ(
	txRep buyticketstxrep.BuyTicketsTxRep,
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	config cnfg.AppConfig,
	authZ auth.AuthZ,
	userRep userrep.UserRep,
	eventRep eventrep.EventRep,
) (BuyTicketsServ, error) {
	return &buyTicketsServ{
		txRep:         txRep,
		tPurchasesRep: tPurchasesRep,
		config:        config,
		authZ:         authZ,
		userRep:       userRep,
		eventRep:      eventRep,
	}, nil
}

func (b *buyTicketsServ) cntFreeTickets(ctx context.Context, eventID uuid.UUID) (int, error) {
	event, err := b.eventRep.GetByID(ctx, eventID)
	if err != nil {
		return 0, fmt.Errorf("checkCntTickets: %v", err)
	}
	txCnt, err := b.txRep.GetCntTxByEventID(ctx, event.GetID())
	if err != nil {
		return 0, fmt.Errorf("checkCntTickets: %v", err)
	}
	purchasesCnt, err := b.tPurchasesRep.GetCntTPurchasesForEvent(ctx, event.GetID())
	if err != nil {
		return 0, fmt.Errorf("checkCntTickets: %v", err)
	}
	freeCnt := event.GetTicketCount() - txCnt - purchasesCnt
	if freeCnt < 0 {
		panic("checkCntTickets: negativ tickets count!?")
	}
	return freeCnt, nil
}

// Если в ctx есть информация об аутентифицированном пользователе то поля customerName, customerEmail не используются
// not sesrver errors: ErrNoFreeTicket, ErrNoUserData, ErrExpireTx
func (b *buyTicketsServ) BuyTicket(
	ctx context.Context,
	eventID uuid.UUID,
	cntTickets int,
	customerName string,
	customerEmail string,
) (*models.TicketPurchaseTx, error) {
	var err error
	ticketsFree, err := b.cntFreeTickets(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}
	if ticketsFree <= 0 {
		return nil, fmt.Errorf("BuyTicket: %w", ErrNoFreeTicket)
	}

	userID, err := b.authZ.UserIDFromContext(ctx)
	if err != nil && err != auth.ErrNotAuthZ {
		return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
	}

	var userName, userEmail string
	if err == auth.ErrNotAuthZ {
		if customerName == "" || customerEmail == "" {
			return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, ErrNoUserData)
		}
		userName = customerName
		userEmail = customerEmail
		userID = uuid.Nil
	} else {
		user, err := b.userRep.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
		}
		userName = user.GetUsername()
		userEmail = user.GetEmail()
		userID = user.GetID()
	}
	timeExpire := time.Now().Add(b.config.BuyTicketTransactionDuration)
	tx, err := models.NewBuyTicketTx(
		uuid.New(),
		userName,
		userEmail,
		time.Now(),
		eventID,
		userID,
		cntTickets,
		timeExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
	}

	err = b.txRep.Add(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}
	return &tx, nil
}

func (b *buyTicketsServ) ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	tx, err := b.txRep.GetByID(ctx, TxID)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}
	ticketPurchase := tx.GetTicketPurchase()
	// --
	// fmt.Printf("TicketPurchase: %+v\n", ticketPurchase)
	// cnt, err := b.txRep.GetCntTxByEventID(ctx, tx.GetTicketPurchase().GetEventID())
	// if err != nil {
	// 	return fmt.Errorf("ConfirmBuyTicket: %v", err)
	// }
	// fmt.Printf("CNT 1: %d\n", cnt)
	// --
	err = b.tPurchasesRep.Add(ctx, ticketPurchase)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}
	err = b.txRep.Delete(ctx, TxID)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}

	// --
	// cnt, err = b.txRep.GetCntTxByEventID(ctx, tx.GetTicketPurchase().GetEventID())
	// if err != nil {
	// 	return fmt.Errorf("ConfirmBuyTicket: %v", err)
	// }
	// fmt.Printf("CNT2: %d\n", cnt)
	// --

	return nil
}

func (b *buyTicketsServ) CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	return b.txRep.Delete(ctx, TxID)
}

func (b *buyTicketsServ) GetAllTicketPurchasesOfUser(
	ctx context.Context,
) ([]*models.TicketPurchase, error) {
	userID, err := b.authZ.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
	}
	tPurchases, err := b.tPurchasesRep.GetTPurchasesOfUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
	}
	return tPurchases, err
}

func (b *buyTicketsServ) GetBuyTicketTransactionDuration() time.Duration {
	return b.config.BuyTicketTransactionDuration
}
