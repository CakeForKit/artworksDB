package buyticketstxrep

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisBuyTicketsTxRep struct {
	rdb *redis.Client
}

var (
	repInstance *RedisBuyTicketsTxRep
	repOnce     sync.Once
)

func NewRedisBuyTicketsTxRep(
	ctx context.Context,
	redisCreds *cnfg.RedisCredentials,
) (*RedisBuyTicketsTxRep, error) {
	var resErr error = nil
	repOnce.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisCreds.Host, redisCreds.Port),
			Password: redisCreds.Password,
			Username: redisCreds.Username,
			DB:       0,
		})
		if err := rdb.Ping(ctx).Err(); err != nil {
			resErr = fmt.Errorf("failed to connect to redis server: %v", err)
			return
		}
		repInstance = &RedisBuyTicketsTxRep{rdb: rdb}
	})

	if resErr != nil {
		return nil, resErr
	}

	return repInstance, nil
}

func (r *RedisBuyTicketsTxRep) GetByID(ctx context.Context, txID uuid.UUID) (*models.TicketPurchaseTx, error) {
	data, err := r.rdb.Get(ctx, txID.String()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("redisRep GetByID: %w", ErrTxNotFound)
		}
		return nil, fmt.Errorf("redisRep GetByID: %v", err)
	}

	var tx models.TicketPurchaseTx
	err = tx.FromJson(data)
	return &tx, err
}

func (r *RedisBuyTicketsTxRep) Add(ctx context.Context, tpTx models.TicketPurchaseTx) error {
	data, err := tpTx.Tojson()
	if err != nil {
		return fmt.Errorf("add ticketTx: %v", err)
	}

	expiration := time.Until(tpTx.GetExpiredAt())
	if expiration <= 0 {
		return fmt.Errorf("add ticketTx: %v", ErrExpireTx)
	}

	err = r.rdb.Set(ctx, tpTx.GetID().String(), data, expiration).Err()
	if err != nil {
		return fmt.Errorf("add tx to redis: %v", err)
	}

	return nil
}

func (r *RedisBuyTicketsTxRep) Delete(ctx context.Context, txID uuid.UUID) error {
	return r.rdb.Del(ctx, txID.String()).Err()
}

// GetCntTxByEventID возвращает количество транзакций для указанного eventID
func (r *RedisBuyTicketsTxRep) GetCntTxByEventID(ctx context.Context, eventID uuid.UUID) (int, error) {
	keys, err := r.rdb.Keys(ctx, "*").Result() // Получаем все ключи из Redis
	if err != nil {
		return 0, fmt.Errorf("CountTxByEventID - get all keys: %v", err)
	}

	count := 0
	for _, key := range keys {
		data, err := r.rdb.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue // Пропускаем если ключ исчез
			}
			return 0, fmt.Errorf("CountTxByEventID: %s: %v", key, err)
		}
		var tx models.TicketPurchaseTx
		if err := tx.FromJson(data); err != nil {
			return 0, fmt.Errorf("CountTxByEventID - fromJson: %s: %v", key, err)
		}
		if tx.GetTicketPurchase().GetEventID() == eventID {
			count++
		}
	}

	return count, nil
}

func (r *RedisBuyTicketsTxRep) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r *RedisBuyTicketsTxRep) Close() {
	r.rdb.Close()
}
