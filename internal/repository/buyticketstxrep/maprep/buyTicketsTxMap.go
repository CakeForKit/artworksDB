package maprep

import (
	"sync"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type BuyTicketsTxMap struct {
	mx    sync.Mutex
	txMap map[uuid.UUID]*models.BuyTicketTx // key - BuyTicketTx.id
}

func NewBuyTicketsTxMap() *BuyTicketsTxMap {
	return &BuyTicketsTxMap{
		txMap: make(map[uuid.UUID]*models.BuyTicketTx),
	}
}

func (b *BuyTicketsTxMap) Get(key uuid.UUID) (*models.BuyTicketTx, bool) {
	b.mx.Lock()
	defer b.mx.Unlock()
	val, ok := b.txMap[key]
	return val, ok
}

func (b *BuyTicketsTxMap) Put(key uuid.UUID, value *models.BuyTicketTx) {
	b.mx.Lock()
	defer b.mx.Unlock()
	b.txMap[key] = value
}

func (b *BuyTicketsTxMap) Delete(key uuid.UUID) {
	b.mx.Lock()
	defer b.mx.Unlock()
	delete(b.txMap, key)
}
