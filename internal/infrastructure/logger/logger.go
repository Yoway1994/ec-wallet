package logger

import (
	"context"
	"ec-wallet/configs"
	"ec-wallet/internal/domain"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func ensureLogDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func NewLogger(config *configs.Config) *zap.Logger {
	cfg := config.Logger
	if err := ensureLogDir(cfg.LogDir); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	// 配置 lumberjack 進行日誌輪替
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%04d-%02d-%02d.log", cfg.LogDir, time.Now().Year(), time.Now().Month(), time.Now().Day()),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	})

	// 根據環境選擇編碼器
	var encodeConfig zapcore.EncoderConfig
	if cfg.Env == "dev" {
		encodeConfig = zap.NewDevelopmentEncoderConfig()
		encodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeConfig = zap.NewProductionEncoderConfig()
		encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	// 設置日誌級別
	level := zap.NewAtomicLevel()
	switch cfg.Level {
	case "debug":
		level.SetLevel(zapcore.DebugLevel)
	case "warn":
		level.SetLevel(zapcore.WarnLevel)
	case "error":
		level.SetLevel(zapcore.ErrorLevel)
	default:
		level.SetLevel(zapcore.InfoLevel)
	}

	// 創建核心
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encodeConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
		level,
	)

	// 啟用 caller 和 stacktrace
	options := []zap.Option{}
	options = append(options, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	zapLogger := zap.New(core, options...)
	_ = zap.ReplaceGlobals(zapLogger)

	return zapLogger
}

func FromGinContext(ctx *gin.Context) *zap.Logger {
	if logger, ok := ctx.Get(string(domain.LoggerKey)); ok {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return zapLogger
		}
	}
	return zap.L()
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(string(domain.LoggerKey)).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}
