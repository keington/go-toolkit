package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/10/25 0:34
 * @file: logger.go
 * @description:
 */

var (
	Logger           *zap.SugaredLogger
	logName          = "default"
	logPathStr       = "logs"  // log path
	logMaxSizeStr    = "100"   // log max size
	logMaxBackupsStr = "30"    // log max backups
	logMaxAgeStr     = "1"     // log max age
	logLevelStr      = "debug" // log level
)

// Level defines the severity of a log message.
type Level uint8

const (
	// DebugLevel logs debug messages.
	DebugLevel Level = iota
	// InfoLevel logs informational messages.
	InfoLevel
	// WarnLevel logs warning messages.
	WarnLevel
	// ErrorLevel logs error messages.
	ErrorLevel
	// FatalLevel logs fatal messages.
	FatalLevel
)

// Config holds logger configuration options.
type Config struct {
	LogName    string
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LogLevel   Level
}

// InitializeLogger initialize the logger
func InitializeLogger(name, path string, maxSize, maxBackups, maxAge int, level string) error {

	logConfig := loggerConfigParse(name, path, maxSize, maxBackups, maxAge, level)

	writeSyncer := getLogWriter(logConfig)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)),
		zapcore.Level(logConfig.LogLevel))

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
	return nil
}

// getEncoder logger encoder
func getEncoder() zapcore.Encoder {

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义时间显示
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 自定义行号显示
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "linenum",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "msg",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      customLevelEncoder, // 大写编码器
		EncodeTime:       customTimeEncoder,  // 自定义时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     customCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(config *Config) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", config.LogPath, config.LogName),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// loggerConfigParse parses logger configuration
func loggerConfigParse(name, path string, maxSize, maxBackups, maxAge int, levelStr string) *Config {

	level := LevelFromString(levelStr)
	logConfig := &Config{
		LogName:    name,
		LogPath:    path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LogLevel:   level,
	}
	return logConfig
}

// LevelFromString converts a string to Level
func LevelFromString(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}
