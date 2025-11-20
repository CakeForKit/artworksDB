package buyticketserv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/buyticketstxrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/ticketpurchasesrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/authz"
	"github.com/google/uuid"
)

var (
	ErrBuyTicketsServ = errors.New("buyTicketsServ")
	ErrNoFreeTicket   = errors.New("no free ticket")
	ErrNoUserData     = errors.New("no info about user (customerName, customerEmail)")
)

type BuyTicketsServ interface {
	BuyTicket(ctx context.Context, eventID uuid.UUID, cntTickets int, customerName string, customerEmail string,
	) (*models.TicketPurchaseTx, error)
	ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error
	CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error
	GetAllTicketPurchasesOfUser(ctx context.Context) ([]*models.TicketPurchase, error)
	GetBuyTicketTransactionDuration() time.Duration
}

type buyTicketsServ struct {
	txRep         buyticketstxrep.BuyTicketsTxRep
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep
	config        cnfg.AppConfig
	authZ         authz.AuthZ
	userRep       userrep.UserRep
	eventRep      eventrep.EventRep
}

func NewBuyTicketsServ(
	txRep buyticketstxrep.BuyTicketsTxRep,
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	config cnfg.AppConfig,
	authZ authz.AuthZ,
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
		return 0, fmt.Errorf("checkCntTickets: %w", err)
	}
	txCnt, err := b.txRep.GetCntTxByEventID(ctx, event.GetID())
	if err != nil {
		return 0, fmt.Errorf("checkCntTickets: %w", err)
	}
	purchasesCnt, err := b.tPurchasesRep.GetCntTPurchasesForEvent(ctx, event.GetID())
	if err != nil {
		return 0, fmt.Errorf("checkCntTickets: %w", err)
	}
	freeCnt := event.GetTicketCount() - txCnt - purchasesCnt
	if freeCnt < 0 {

		return 0, fmt.Errorf("checkCntTickets: %w", ErrNoFreeTicket)
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
	if err := b.validateTicketAvailability(ctx, eventID); err != nil {
		return nil, err
	}

	userID, userName, userEmail, err := b.getUserData(ctx, customerName, customerEmail)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuyTicketsServ, err)
	}

	return b.createPurchaseTransaction(ctx, eventID, cntTickets, userID, userName, userEmail)
}

func (b *buyTicketsServ) validateTicketAvailability(ctx context.Context, eventID uuid.UUID) error {
	ticketsFree, err := b.cntFreeTickets(ctx, eventID)
	if err != nil {
		return fmt.Errorf("BuyTicket: %w", err)
	}
	if ticketsFree <= 0 {
		return fmt.Errorf("BuyTicket: %w", ErrNoFreeTicket)
	}
	return nil
}

func (b *buyTicketsServ) getUserData(ctx context.Context, customerName,
	customerEmail string) (uuid.UUID, string, string, error) {
	userID, err := b.authZ.UserIDFromContext(ctx)
	if err != nil && err != authz.ErrNotAuthZ {
		return uuid.Nil, "", "", err
	}

	if err == authz.ErrNotAuthZ {
		return b.getUnauthenticatedUserData(customerName, customerEmail)
	}
	return b.getAuthenticatedUserData(ctx, userID)
}

func (b *buyTicketsServ) getUnauthenticatedUserData(customerName,
	customerEmail string) (uuid.UUID, string, string, error) {
	if customerName == "" || customerEmail == "" {
		return uuid.Nil, "", "", ErrNoUserData
	}
	return uuid.Nil, customerName, customerEmail, nil
}

func (b *buyTicketsServ) getAuthenticatedUserData(ctx context.Context,
	userID uuid.UUID) (uuid.UUID, string, string, error) {
	user, err := b.userRep.GetByID(ctx, userID)
	if err != nil {
		return uuid.Nil, "", "", err
	}
	return user.GetID(), user.GetUsername(), user.GetEmail(), nil
}

func (b *buyTicketsServ) createPurchaseTransaction(
	ctx context.Context,
	eventID uuid.UUID,
	cntTickets int,
	userID uuid.UUID,
	userName string,
	userEmail string,
) (*models.TicketPurchaseTx, error) {
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

	if err := b.txRep.Add(ctx, tx); err != nil {
		return nil, fmt.Errorf("BuyTicket: %w", err)
	}

	return &tx, nil
}

func (b *buyTicketsServ) ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	tx, err := b.txRep.GetByID(ctx, TxID)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %w", err)
	}
	ticketPurchase := tx.GetTicketPurchase()

	err = b.tPurchasesRep.Add(ctx, ticketPurchase)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %w", err)
	}
	err = b.txRep.Delete(ctx, TxID)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %w", err)
	}
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
