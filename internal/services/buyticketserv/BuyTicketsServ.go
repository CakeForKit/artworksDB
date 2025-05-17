package buyticketserv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"github.com/google/uuid"
)

var (
	ErrNoFreeTicket = errors.New("no free ticket")
)

type BuyTicketsServ interface {
	BuyTicket(ctx context.Context, event models.Event, cntTickets int, customerName string, customerEmail string) (*models.TicketPurchaseTx, error)
	BuyTicketByUser(ctx context.Context, event models.Event, cntTickets int, user models.User) (*models.TicketPurchaseTx, error)
	ConfirmBuyTicket(ctx context.Context, TxID uuid.UUID) error
	CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error
	GetAllTicketPurchasesOfUser(ctx context.Context, userID uuid.UUID) ([]*models.TicketPurchase, error)
}

type buyTicketsServ struct {
	txRep         buyticketstxrep.BuyTicketsTxRep
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep
	config        cnfg.AppConfig
}

func NewBuyTicketsServ(
	txRep buyticketstxrep.BuyTicketsTxRep,
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	config cnfg.AppConfig,
) (BuyTicketsServ, error) {
	return &buyTicketsServ{
		txRep:         txRep,
		tPurchasesRep: tPurchasesRep,
		config:        config,
	}, nil
}

func (b *buyTicketsServ) cntFreeTickets(ctx context.Context, event models.Event) (int, error) {
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

func (b *buyTicketsServ) BuyTicket(
	ctx context.Context,
	event models.Event,
	cntTickets int,
	customerName string,
	customerEmail string,
) (*models.TicketPurchaseTx, error) {
	ticketsFree, err := b.cntFreeTickets(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}
	if ticketsFree <= 0 {
		return nil, fmt.Errorf("BuyTicket: %v", ErrNoFreeTicket)
	}

	timeExpire := time.Now().Add(b.config.BuyTicketTransactionDuration)
	tx, err := models.NewBuyTicketTx(
		uuid.New(),
		customerName,
		customerEmail,
		time.Now(),
		event.GetID(),
		uuid.Nil,
		cntTickets,
		timeExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}

	// // ----
	// fmt.Printf("TX         : %+v\n", tx)
	// j, err := tx.Tojson()
	// if err != nil {
	// 	return nil, fmt.Errorf("BuyTicket: %v", err)
	// }
	// var outtx models.TicketPurchaseTx
	// err = outtx.FromJson(j)
	// if err != nil {
	// 	return nil, fmt.Errorf("BuyTicket: %v", err)
	// }
	// fmt.Printf("TX fromJson: %+v\n", outtx)
	// fmt.Printf("tx   : %s\n", tx.GetExpiredAt().String())
	// fmt.Printf("outtx: %s\n", outtx.GetExpiredAt().String())
	// // ----

	err = b.txRep.Add(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}
	return &tx, nil
}

func (b *buyTicketsServ) BuyTicketByUser(
	ctx context.Context,
	event models.Event,
	cntTickets int,
	user models.User,
) (*models.TicketPurchaseTx, error) {
	ticketsFree, err := b.cntFreeTickets(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
	}
	if ticketsFree <= 0 {
		return nil, fmt.Errorf("BuyTicket: %v", ErrNoFreeTicket)
	}

	timeExpire := time.Now().Add(b.config.BuyTicketTransactionDuration)
	tx, err := models.NewBuyTicketTx(
		uuid.New(),
		user.GetUsername(),
		user.GetEmail(),
		time.Now(),
		event.GetID(),
		user.GetID(),
		cntTickets,
		timeExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("BuyTicket: %v", err)
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
	fmt.Printf("TicketPurchase: %+v\n", ticketPurchase)
	cnt, err := b.txRep.GetCntTxByEventID(ctx, tx.GetTicketPurchase().GetEventID())
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}
	fmt.Printf("CNT 1: %d\n", cnt)
	// --
	err = b.txRep.Delete(ctx, TxID)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}

	// --
	cnt, err = b.txRep.GetCntTxByEventID(ctx, tx.GetTicketPurchase().GetEventID())
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}
	fmt.Printf("CNT2: %d\n", cnt)
	// --

	err = b.tPurchasesRep.Add(ctx, ticketPurchase)
	if err != nil {
		return fmt.Errorf("ConfirmBuyTicket: %v", err)
	}

	return nil
}

func (b *buyTicketsServ) CancelBuyTicket(ctx context.Context, TxID uuid.UUID) error {
	return b.txRep.Delete(ctx, TxID)
}

func (b *buyTicketsServ) GetAllTicketPurchasesOfUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*models.TicketPurchase, error) {
	tPurchases, err := b.tPurchasesRep.GetTPurchasesOfUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAllTicketPurchasesOfUser: %v", err)
	}
	return tPurchases, err
}
