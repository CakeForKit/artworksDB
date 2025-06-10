package cnfg

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// Viper использует пакет mapstructure под капотом для преобразования значений

const (
	PostgresDB   = "postgres"
	ClickHouseDB = "clickhouse"
)

type AppConfig struct {
	Datebase                     string        `mapstructure:"datebase"`
	TokenSymmetricKey            string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration          time.Duration `mapstructure:"access_token_duration"`
	BuyTicketTransactionDuration time.Duration `mapstructure:"buy_ticket_transaction_duration"`
	Port                         int           `mapstructure:"port"`
}

type DatebaseConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type DatebaseCredentials struct {
	Host     string
	DbName   string
	Port     int
	Username string
	Password string
}

type PostgresTestConfig struct {
	DbName       string `mapstructure:"postgres_db"`
	Port         int    `mapstructure:"postgres_port"`
	Username     string `mapstructure:"postgres_user"`
	Password     string `mapstructure:"postgres_password"`
	Image        string `mapstructure:"postgres_image"`
	MigrationDir string `mapstructure:"postgres_migration_dir"`
}

type RedisCredentials struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Username string `mapstructure:"REDIS_USER"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

var (
	ErrConfigRead    = errors.New("ReadInConfig")
	ErrUnmarshalRead = errors.New("err to unmarshal config ")
	ErrEnvRead       = errors.New("read env error")
	ErrUnknownDB     = errors.New("unknown datebase")
)

func LoadAppConfig() (config *AppConfig, err error) {
	config = &AppConfig{}
	v := viper.New()
	v.AddConfigPath("./configs/")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = v.UnmarshalKey("app", config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if config.Datebase != PostgresDB && config.Datebase != ClickHouseDB {
		return nil, ErrUnknownDB
	}

	return config, nil
}

func LoadDatebaseConfig(path string) (config *DatebaseConfig, err error) {
	config = &DatebaseConfig{}
	v := viper.New()
	v.AddConfigPath(path) // "./configs/"
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = v.UnmarshalKey("datebase", config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	return config, nil
}

type PostgresCredentials struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	DbName   string `mapstructure:"POSTGRES_DB"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
}

func LoadPgCredentials(path string) (*DatebaseCredentials, error) {
	viper.AddConfigPath(path) // расположение файла с настройками
	viper.SetConfigName("db")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	config := &PostgresCredentials{}
	if err := viper.Unmarshal(config); err != nil { // преобразование значений в переданный объект
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	return &DatebaseCredentials{
		Host:     config.Host,
		DbName:   config.DbName,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}, nil
}

func LoadPgTestConfig() (config *PostgresTestConfig, err error) {
	config = &PostgresTestConfig{}
	v := viper.New()
	v.AddConfigPath("./configs/") // расположение файла с настройками
	v.SetConfigName("db_test_config")
	v.SetConfigType("yaml")
	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	if err = v.UnmarshalKey("postgres", config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	fmt.Printf("LoadPgTestConfig: %+v\n", config)

	return config, nil
}

func LoadRedisCredentials() (config *RedisCredentials, err error) {
	viper.AddConfigPath("./configs/") // расположение файла с настройками
	viper.SetConfigName("redis")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	config = &RedisCredentials{}
	if err = viper.Unmarshal(config); err != nil { // преобразование значений в переданный объект
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	return config, nil
}

type ClickHouseCredentials struct {
	Host     string `mapstructure:"CLICKHOUSE_HOST"`
	DbName   string `mapstructure:"CLICKHOUSE_DB"`
	Port     int    `mapstructure:"CLICKHOUSE_PORT"`
	Username string `mapstructure:"CLICKHOUSE_USER"`
	Password string `mapstructure:"CLICKHOUSE_PASSWORD"`
}

func LoadClickHouseCredentials() (*DatebaseCredentials, error) {
	viper.AddConfigPath("./configs/") // расположение файла с настройками
	viper.SetConfigName("clickhouse")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	config := &ClickHouseCredentials{}
	if err := viper.Unmarshal(config); err != nil { // преобразование значений в переданный объект
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	return &DatebaseCredentials{
		Host:     config.Host,
		DbName:   config.DbName,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}, nil
}

func GetProjectRoot() string {
	_, currentFile, _, _ := runtime.Caller(0) // Получаем путь к текущему файлу
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	return projectRoot
}

func GetPgTestConfig() (config *PostgresTestConfig) {
	projectRoot := GetProjectRoot()
	migrationDir := filepath.Join(projectRoot, "migrations") // Путь от корня проекта
	return &PostgresTestConfig{
		DbName:       "testartwork",
		Port:         5432,
		Username:     "testUser",
		Password:     "testPassword",
		Image:        "postgres:latest",
		MigrationDir: migrationDir,
	}
}

func GetTestDatebaseConfig() (config *DatebaseConfig) {
	return &DatebaseConfig{
		MaxOpenConns:    2,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Hour,
	}
}
