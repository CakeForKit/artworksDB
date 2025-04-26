package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Viper использует пакет mapstructure под капотом для преобразования значений

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Datebase DatebaseConfig
}

type AppConfig struct {
	TokenSymmetricKey            string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration          time.Duration `mapstructure:"access_token_duration"`
	BuyTicketTransactionDuration time.Duration `mapstructure:"buy_ticket_transaction_duration"`
}

type DatebaseConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	DbName   string `mapstructure:"POSTGRES_DB"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
}

var (
	ErrConfigRead = errors.New("ReadInConfig")
	ErrEnvRead    = errors.New("read env error")
)

func LoadConfig() (config *Config, err error) {
	config = &Config{}
	viper.AddConfigPath("./configs/") // расположение файла с настройками
	viper.SetConfigName("db")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	cnfPostgres := &PostgresConfig{}
	if err = viper.Unmarshal(cnfPostgres); err != nil { // преобразование значений в переданный объект
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	config.Postgres = *cnfPostgres

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	return config, nil
}
