package logger

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)

	return r.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {

	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	return func(c *gin.Context) {
		start := time.Now()
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		req := c.Request

		user, ok := c.Get("User")
		if !ok {
			user = nil
		}

		fields := []zapcore.Field{
			zap.String("time", time.Now().String()),
			zap.Int("status", w.Status()),
			zap.String("latency", time.Since(start).String()),
			zap.Any("user", user),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("response", w.body.String()),
			zap.String("inFunction", c.HandlerName()),
			zap.String("host", req.Host),
		}

		n := w.Status()

		switch {
		case n >= 500:
			logger.Error("Server error", fields...)
		case n >= 400:
			logger.Warn("Client error", fields...)
		case n >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Success", fields...)
		}

	}
}
