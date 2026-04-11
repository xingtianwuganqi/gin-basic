package logger

import (
	"gin-basic/settings"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	encoder := getEncoder()
	logLevel := getLogLevel(settings.Conf.App.LogLevel)

	stdoutLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= logLevel && level < zapcore.ErrorLevel
	})
	stderrLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= logLevel && level >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), stdoutLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), stderrLevel),
	)

	Logger = zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed.
func Debug(message string, fields ...zap.Field) {
	Logger.Debug(message, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed.
func Info(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed.
func Warn(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed.
func Error(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed.
func Panic(message string, fields ...zap.Field) {
	Logger.Panic(message, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed.
func Fatal(message string, fields ...zap.Field) {
	Logger.Fatal(message, fields...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	if Logger == nil {
		return nil
	}
	return Logger.Sync()
}
