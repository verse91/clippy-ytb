package logger

import (
	"os"

	// z "go.uber.org/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// go run cmd/cli/main.log.go
func InitLogger() {
	// sugar := z.NewExample().Sugar()
	// sugar.Infof("Hello: %s, Age: %d", "Verse", 100)

	// logger := z.NewExample()
	// logger.Info("Hello", z.String("name", "Verse"), z.Int("age", 100))

	// logger := z.NewExample()
	// // logger.Info("Hello")

	// // // Development
	// // logger, _ = z.NewDevelopment()
	// // logger.Info("This is Development")

	// // Production
	// logger, _ = z.NewProduction()
	// logger.Info("This is Production")

	encoder := getEncoderLog()
	writerSync := getWriterSync()
	core := zapcore.NewCore(encoder, writerSync, zapcore.InfoLevel)
	Log = zap.New(core, zap.AddCaller())

	Log.Info("Info log", zap.Int("line", 1))
	Log.Error("Error log", zap.Int("line", 2))

}

// format logs a msg
func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func getWriterSync() zapcore.WriteSyncer {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll("./pkg/logger/log", 0755); err != nil {
		panic(err)
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./pkg/logger/log/info.log",
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5,   //days
		Compress:   true, // disabled by default
        LocalTime: true,
	}
	syncfile := zapcore.AddSync(lumberjackLogger)
	syncConsole := zapcore.AddSync(os.Stderr)
	return zapcore.NewMultiWriteSyncer(syncConsole, syncfile)

}
