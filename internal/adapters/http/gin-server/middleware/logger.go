package middleware

import (
	"ec-wallet/internal/adapters/http/gin-server/utils"
	"ec-wallet/internal/domain"
	"ec-wallet/internal/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDMiddleware 為每個請求添加唯一的 request ID
func RequestIDMiddleware(headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(headerName)
		if requestID == "" {
			uuid, err := uuid.NewRandom()
			if err != nil {
				utils.HandleError(c, errors.ErrRequestIDGeneration)
				return
			}
			requestID = uuid.String()
		}
		c.Set(string(domain.RequestIDKey), requestID)
		c.Header(headerName, requestID)
		c.Next()
	}
}

// LoggerConfig defines configuration for LoggerMiddleware
type LoggerConfig struct {
	LogQueryParams    bool
	LogUserAgent      bool
	LogRequestURLPath bool
	LogClientIP       bool
	LogMethod         bool
}

// LoggerMiddleware 為每個請求添加帶有 request ID 的 logger
func LoggerMiddleware(baseLogger *zap.Logger, config *LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, exists := c.Get(string(domain.RequestIDKey))
		requestIDStr, ok := requestID.(string)
		if !exists || !ok {
			requestIDStr = "unknown"
			baseLogger.Error("Invalid or missing request_id in context")
		}
		//
		fields := []zap.Field{
			zap.String("request_id", requestIDStr),
		}
		if config != nil {
			if config.LogClientIP {
				fields = append(fields, zap.String("client_ip", c.ClientIP()))
			}
			if config.LogRequestURLPath {
				fields = append(fields, zap.String("path", c.Request.URL.Path))
			}
			if config.LogQueryParams && c.Request.URL.RawQuery != "" {
				fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
			}
			if config.LogUserAgent {
				fields = append(fields, zap.String("user_agent", c.Request.UserAgent()))
			}
			if config.LogMethod {
				fields = append(fields, zap.String("method", c.Request.Method))
			}
		}

		reqLogger := baseLogger.With(fields...)
		c.Set(string(domain.LoggerKey), reqLogger)
		//
		reqLogger.Info("Request started")
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		status := c.Writer.Status()
		logFields := []zap.Field{
			zap.Int("status", status),
			zap.Duration("latency", latency),
		}
		if status >= 500 {
			reqLogger.Error("Request completed with server error", logFields...)
		} else if status >= 400 {
			if len(c.Errors) > 0 {
				reqLogger.Warn("Request failed", append(logFields, zap.Any("errors", c.Errors))...)
			} else {
				reqLogger.Warn("Request completed with client error", logFields...)
			}
		} else {
			reqLogger.Info("Request completed", logFields...)
		}
		reqLogger.Sync()
	}
}
