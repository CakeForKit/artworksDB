package testobj

import (
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
)

type AppConfigMother interface {
	Default() cnfg.AppConfig
}

func NewAppConfigMother() AppConfigMother {
	return &appConfigMother{}
}

type appConfigMother struct{}

func (am *appConfigMother) Default() cnfg.AppConfig {
	config := cnfg.AppConfig{
		Datebase:                     cnfg.PostgresDB,
		TokenSymmetricKey:            "01234567890123456789012345678912",
		AccessTokenDuration:          time.Hour,
		BuyTicketTransactionDuration: time.Hour,
		Port:                         8000,
	}
	return config
}
