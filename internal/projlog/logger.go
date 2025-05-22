package projlog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrDirLog        = errors.New("error create dir logs")
	ErrTypeLog       = errors.New("error type log")
	ErrLevelLog      = errors.New("error level log")
	ErrCreateEncoder = errors.New("error create encoder")
	ErrBuildLogger   = errors.New("error build logger")
	ErrSugarLogger   = errors.New("error sugar logger")
)

const (
	DevTypeLog  = "dev"
	ProdTypeLog = "prod"
)

const (
	InfoLevelLog  = "info"
	DebugLevelLog = "debug"
	WarnLevelLog  = "warn"
	ErrorLevelLog = "error"
	PanicLevelLog = "panic"
	FatalLevelLog = "fatal"
)

func getEncoderConf(typeStr string) (*zapcore.EncoderConfig, error) {
	var encoderCfg zapcore.EncoderConfig
	switch typeStr {
	case ProdTypeLog:
		encoderCfg = zap.NewProductionEncoderConfig()
		// encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		// encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	case DevTypeLog:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		// encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		// }
	default:
		return nil, ErrTypeLog
	}
	encoderCfg.TimeKey = "time"
	encoderCfg.MessageKey = "msg"
	encoderCfg.CallerKey = "caller"
	return &encoderCfg, nil
}

func getLevelLog(levelStr string) (*zapcore.Level, error) {
	// Определяем уровень логирования из конфига
	var level zapcore.Level
	switch levelStr {
	case DebugLevelLog:
		level = zapcore.DebugLevel
	case InfoLevelLog:
		level = zapcore.InfoLevel
	case WarnLevelLog:
		level = zapcore.WarnLevel
	case ErrorLevelLog:
		level = zapcore.ErrorLevel
	case PanicLevelLog:
		level = zapcore.PanicLevel
	case FatalLevelLog:
		level = zapcore.FatalLevel
	default:
		return nil, ErrLevelLog
	}
	return &level, nil
}

func NewLogger(logCnfg *cnfg.LogConfig) (*zap.SugaredLogger, error) {
	// Создаём директорию для логов, если её нет
	logDir := filepath.Dir(logCnfg.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("NewLogger: %w", ErrDirLog)
	}

	encoderCfg, err := getEncoderConf(logCnfg.Type)
	if err != nil {
		return nil, fmt.Errorf("NewLogger: %w", err)
	}

	level, err := getLevelLog(logCnfg.Level)
	if err != nil {
		return nil, fmt.Errorf("NewLogger: %w", err)
	}

	fileEncoder := zapcore.NewJSONEncoder(*encoderCfg)
	if fileEncoder == nil {
		return nil, fmt.Errorf("file encoder: %w", ErrCreateEncoder)
	}

	consoleEncoderCfg := *encoderCfg
	consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	consoleEncoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	consoleEncoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)
	if consoleEncoder == nil {
		return nil, fmt.Errorf("console encoder: %w", ErrCreateEncoder)
	}

	// Ядра для логирования (файл + консоль)
	cores := []zapcore.Core{}

	// Настройка вывода в файл (текстовый формат)
	logFile, _ := os.OpenFile(
		logCnfg.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	fileWriter := zapcore.AddSync(logFile)

	cores = append(cores, zapcore.NewCore(
		fileEncoder,
		fileWriter,
		level,
	))

	if logCnfg.Stdout {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// Создаём мультиплексор
	core := zapcore.NewTee(cores...)

	// Собираем логгер
	logger := zap.New(
		core,
		zap.AddCaller(),                       // Добавляем caller (откуда вызван лог)
		zap.AddStacktrace(zapcore.ErrorLevel), // Трассировка для ошибок
	)
	if logger == nil {
		return nil, fmt.Errorf("NewLogger: %w", ErrBuildLogger)
	}

	sugarLog := logger.Sugar()
	if sugarLog == nil {
		return nil, fmt.Errorf("NewLogger: %w", ErrSugarLogger)
	}
	return sugarLog, nil
}
