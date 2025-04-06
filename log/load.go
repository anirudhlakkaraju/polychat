package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/anirudhlakkaraju/polychat/config/props"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

var (
	onceInit          = new(sync.Once)
	implementation    = make(map[string]interface{})
	customLoggerKey   = "customLogger"
	attachedLoggerKey = "customContextLogger"
)

func Init() error {
	var err error
	onceInit.Do(func() {
		if implementation[customLoggerKey] == nil {
			props := props.GetProps()
			handler := customHandler(props)
			implementation[customLoggerKey] = slog.New(handler)
		}
	})
	return err
}

func customHandler(props *properties.Properties) slog.Handler {

	var level slog.Level
	switch props.GetString("log.level", "info") {
	case "debug":
		level = slog.LevelDebug
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: props.GetBool("log.source", false),
	}

	if path := props.GetString("log.file.path", ""); path != "" {
		// Open file and create a file handler with JSON format
		logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		// Create a multi writer to log to both console and logFile
		multiWriter := io.MultiWriter(os.Stdout, logFile)

		return slog.NewJSONHandler(multiWriter, opts)
	}

	return slog.NewJSONHandler(os.Stdout, opts)
}

// GetCustomLogger returns the configured logger
func GetCustomLogger() *slog.Logger {
	v := implementation[customLoggerKey]
	return v.(*slog.Logger)
}

// AttachLoggerToContext ties the custom logger to gin context
func AttachLoggerToContext(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Add request ID field to logger

		ctxWithLogger := context.WithValue(ctx, attachedLoggerKey, logger)
		ctx.Request = ctx.Request.WithContext(ctxWithLogger)
		ctx.Next()
	}
}

// LoggerFromContext returns the logger that's tied to the context
// If not it returns the custom logger
func LoggerFromContext(ctx context.Context) *slog.Logger {

	defaultLogger := GetCustomLogger()

	if ctx == nil {
		defaultLogger.Error("context is nil, defaulting to custom logger")
		return defaultLogger
	}

	logger, ok := ctx.Value(attachedLoggerKey).(*slog.Logger)
	if ok && logger != nil {
		return logger
	}

	return defaultLogger
}
