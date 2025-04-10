package util

import (
	"time"

	"github.com/spf13/viper"
)

// Viper использует пакет mapstructure под капотом для преобразования значений

type Config struct {
	TokenSymmetricKey            string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration          time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration         time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	BuyTicketTransactionDuration time.Duration `mapstructure:"BUY_TICKET_TRANSACTION_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // расположение файла с настройками
	viper.SetConfigName("app") // искать файл с определенным именем
	viper.SetConfigType("env") // тип файла
	viper.AutomaticEnv()       // чтобы Viper автоматически переопределял значения, которые он прочитал из файла с настройками, значениями соответствующих переменных окружения, если они существуют.

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // преобразование значений в переданный объект
	return
}
