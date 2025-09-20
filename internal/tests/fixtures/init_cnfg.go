package fixtures

import (
	"errors"
	"fmt"
	"sync"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
)

var (
	ErrInitTestConfig = errors.New("InitTestConfig")
)

var (
	initOnce   sync.Once
	appCnfg    *cnfg.AppConfig
	dbCreds    *cnfg.DatebaseCredentials
	dbCnfg     *cnfg.DatebaseConfig
	redisCreds *cnfg.RedisCredentials
	initError  error
)

// InitTestConfig инициализирует конфигурацию один раз для всех тестов
func InitTestConfig(dirConfig string) (*cnfg.AppConfig, *cnfg.DatebaseCredentials, *cnfg.DatebaseConfig, *cnfg.RedisCredentials, error) {
	initOnce.Do(func() {
		appCnfg, initError = cnfg.LoadAppConfig(dirConfig, "config", "yaml")
		if initError != nil {
			initError = fmt.Errorf("%w: %w", ErrInitTestConfig, initError)
			return
		}

		switch appCnfg.Datebase {
		case cnfg.PostgresDB:
			pgCreds, err := cnfg.LoadPgCredentials(dirConfig, "test_db", "env")
			if err != nil {
				initError = fmt.Errorf("%w: %w", ErrInitTestConfig, err)
				return
			}
			dbCreds = pgCreds
		case cnfg.ClickHouseDB:
			clhCreds, err := cnfg.LoadClickHouseCredentials(dirConfig, "clickhouse", "env")
			if err != nil {
				initError = fmt.Errorf("%w: %w", ErrInitTestConfig, err)
				return
			}
			dbCreds = clhCreds
		}

		redisCreds, initError = cnfg.LoadRedisCredentials(dirConfig, "redis", "env")
		if initError != nil {
			initError = fmt.Errorf("%w: %w", ErrInitTestConfig, initError)
			return
		}

		dbCnfg, initError = cnfg.LoadDatebaseConfig(dirConfig, "config", "yaml")
		if initError != nil {
			initError = fmt.Errorf("%w: %w", ErrInitTestConfig, initError)
			return
		}
	})

	return appCnfg, dbCreds, dbCnfg, redisCreds, initError
}

// func RunInTx(t *testing.T, fn func(tx *sql.Tx)) {
// 	t.Helper()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		t.Fatalf("Не удалось начать транзакцию: %v", err)
// 	}
// 	defer func() { _ = tx.Rollback() }()
// 	fn(tx)
// }
