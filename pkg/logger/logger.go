package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugarLogger *zap.SugaredLogger

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[viper.GetString("logger.level")]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}

func InitLogger() {
	if sugarLogger != nil {
		return
	}

	logFile := viper.GetString("logger.logfile")
	logLevel := getLoggerLevel()

	var logWriter zapcore.WriteSyncer

	if logFile != "stdout" {
		logWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:  logFile,
			MaxSize:   1 << 30, //1G
			LocalTime: true,
			Compress:  true,
		})
	} else {
		logWriter = zapcore.AddSync(os.Stdout)
	}

	var encoder zapcore.EncoderConfig

	if logLevel == zapcore.DebugLevel {
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder.EncodeTime = syslogTimeEncoder
		encoder.EncodeCaller = debugCaller
		encoder.StacktraceKey = "stack"
	} else {
		encoder = zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = logger.Sugar()
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan _2 15:04:05.000000"))
}

func debugCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(filepath.Base(caller.FullPath()) + " " + strconv.Itoa(os.Getpid()))
}

func Debug(args ...interface{}) {
	sugarLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	sugarLogger.With("bar", "style").Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	sugarLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugarLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	sugarLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	sugarLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	sugarLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	sugarLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	sugarLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}
