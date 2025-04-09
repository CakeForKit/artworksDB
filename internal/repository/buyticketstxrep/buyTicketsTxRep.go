package buyticketstxrep

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep/maprep"
	"github.com/google/uuid"
)

type BuyTicketsTxRep interface {
	Get(key uuid.UUID) (*models.BuyTicketTx, bool)
	Put(key uuid.UUID, value *models.BuyTicketTx)
	Delete(key uuid.UUID)
}

func NewBuyTicketsTxRep() (BuyTicketsTxRep, error) {
	return maprep.NewBuyTicketsTxMap(), nil
}
