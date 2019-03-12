package logger

import (
	"github.com/Baldomo/Fangs/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var jack = &lumberjack.Logger{
	Filename:   "/var/log/fangs/fangs.log",
	MaxSize:    100,
	MaxBackups: 20,
	MaxAge:     30,
	Compress:   true,
	LocalTime:  true,
}

var l *zap.SugaredLogger

func init() {
	var logger *zap.Logger

	if utils.IsDebug() {
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)

		logger = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	} else {
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(jack),
			zap.InfoLevel,
		)

		logger = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	}

	l = logger.Sugar()
}

func Debug(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Debugw(message, keysAndValues...)
}

func DPanic(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.DPanicw(message, keysAndValues...)
}

func Error(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Errorw(message, keysAndValues...)
}

func Fatal(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Fatalw(message, keysAndValues...)
}

func Info(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Infow(message, keysAndValues...)
}

func Panic(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Panicw(message, keysAndValues...)
}

func Warn(message string, keysAndValues ...interface{}) {
	defer l.Sync()
	l.Warnw(message, keysAndValues...)
}
