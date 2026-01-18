package logger

import (
	"gin-basic/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
	"os"
)

var Logger *zap.Logger

func InitLogger() {
	// 获取配置
	logConfig := settings.Conf.Log

	// 创建Encoder
	encoder := getEncoder()

	// 创建WriteSyncer
	writeSyncer := getLogWriter(logConfig.Filename, logConfig.MaxSize, logConfig.MaxBackups, logConfig.MaxAge, logConfig.Compress)

	// 创建Core
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

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

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return zapcore.AddSync(lumberJackLogger)
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
	return Logger.Sync()
}