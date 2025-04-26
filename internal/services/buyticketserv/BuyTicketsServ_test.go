package buyticketserv

import (
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep/maprep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuyTicketsService(t *testing.T) {
	validUserID := uuid.New()
	validEventID := uuid.New()
	validCntTickets := 2
	invalidCntTickets := 0
	validTxID := uuid.New()

	t.Run("BuyTicket - valid", func(t *testing.T) {
		txRep := maprep.NewBuyTicketsTxMap()

		config := config.Config{
			App: config.AppConfig{
				BuyTicketTransactionDuration: 10 * time.Minute,
			},
		}

		serv := &buyTicketsServ{
			txRep:  txRep,
			config: config,
		}

		err := serv.BuyTicket(validUserID, validEventID, validCntTickets)
		require.NoError(t, err)
	})

	t.Run("BuyTicket - invalid cntTickets", func(t *testing.T) {
		serv := &buyTicketsServ{}
		err := serv.BuyTicket(validUserID, validEventID, invalidCntTickets)
		require.Error(t, err)
		require.Equal(t, models.ErrBuyTicketTxZeroCnt, err)
	})

	t.Run("ConfirmBuyTicket", func(t *testing.T) {
		txRep := maprep.NewBuyTicketsTxMap()

		serv := &buyTicketsServ{
			txRep: txRep,
		}

		serv.ConfirmBuyTicket(validTxID)
	})
}
