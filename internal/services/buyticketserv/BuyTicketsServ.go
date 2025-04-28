package buyticketserv

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"github.com/google/uuid"
)

// поиск для подтверждения покупки
// поиск для удаления завершившихся по времени транзакций

type BuyTicketsServ interface {
	BuyTicket(userID uuid.UUID, eventID uuid.UUID, cntTickets int) error
	ConfirmBuyTicket(TxID uuid.UUID)
}

func NewBuyTicketsServ(txRep buyticketstxrep.BuyTicketsTxRep, config cnfg.AppConfig) (BuyTicketsServ, error) {
	return &buyTicketsServ{
		txRep:  txRep,
		config: config,
	}, nil
}

type buyTicketsServ struct {
	txRep  buyticketstxrep.BuyTicketsTxRep
	config cnfg.AppConfig
}

func (b *buyTicketsServ) BuyTicket(userID uuid.UUID, eventID uuid.UUID, cntTickets int) error {
	timeExpire := time.Now().Add(b.config.BuyTicketTransactionDuration)
	tx, err := models.NewBuyTicketTx(uuid.New(), userID, eventID, cntTickets, timeExpire)
	if err != nil {
		return err
	}
	b.txRep.Put(tx.GetID(), &tx)
	return nil
}

func (b *buyTicketsServ) ConfirmBuyTicket(TxID uuid.UUID) {
	b.txRep.Delete(TxID)
}
