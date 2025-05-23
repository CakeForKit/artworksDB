package cnfg

import (
	"fmt"

	"github.com/spf13/viper"
)

type LogConfig struct {
	Type   string `mapstructure:"type"`
	Level  string `mapstructure:"level"`
	File   string `mapstructure:"file"`
	Stdout bool   `mapstructure:"stdout"`
}

// func GetTestLogConfig() *LogConfig {
// 	return &LogConfig{
// 		Type:   "dev",
// 		Level:  "debug",
// 		File:   "logs/app.log",
// 		Stdout: true,
// 	}
// }

func GetLogConfig() (*LogConfig, error) {
	config := &LogConfig{}
	v := viper.New()
	v.AddConfigPath("./configs/") // Папка с конфигами
	v.SetConfigName("logger")     // Имя файла без расширения
	v.SetConfigType("yaml")       // Тип файла

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalRead, err)
	}

	return config, nil
}
