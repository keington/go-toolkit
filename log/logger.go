package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/10/24 22:15
 * @file: log.go
 * @description:
 */

// Logger is the logger instance.
var logger *zap.SugaredLogger

// LogConfig holds log configuration options.
type LogConfig struct {
	LogName    string
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LogLevel   zapcore.Level
}

// InitializeLogger initializes the log with the given configuration.
func InitializeLogger(config LogConfig) error {
	writeSyncer := getLogWriter(&config)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
		config.LogLevel)

	log := zap.New(core, zap.AddCaller())
	logger = log.Sugar()
	return nil
}

// NewLogConfig creates a LogConfig with default values.
func NewLogConfig() LogConfig {
	return LogConfig{
		LogName:    "default",
		LogPath:    "logs",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     1,
		LogLevel:   zapcore.DebugLevel,
	}
}

// logger encoder
func getEncoder() zapcore.Encoder {

	// 日志级别显示
	setLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 时间显示
	setTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 行号显示
	setCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "log",
		CallerKey:        "caller",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "msg",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      setLevelEncoder, // 大写编码器
		EncodeTime:       setTimeEncoder,  // 自定义时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     setCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(config *LogConfig) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", config.LogPath, config.LogName),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
