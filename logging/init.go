package logging

import (
	"github.com/aka-achu/eidos"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

// Declaring logger variable for global access
var (
	AppLogger     *zap.SugaredLogger
	RequestLogger *zap.SugaredLogger
	RepoLogger    *zap.SugaredLogger
)

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	l, err := eidos.New(fileName, &eidos.Options{
		Size:             10,
		Period:           0,
		RetentionPeriod:  30,
		Compress:         true,
		CompressionLevel: 0,
		LocalTime:        true,
	}, &eidos.Callback{})
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(l)
}

func Initialize() {

	path, _ := os.Getwd()
	path = filepath.Join(path, "log")

	// Validating the existence of the logging path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal("Failed to create the folder for storing the logs", err)
		}
	}

	// Initializing the declared logging variables
	AppLogger = zap.New(
		zapcore.NewCore(
			getEncoder(),
			getLogWriter(filepath.Join(path, "app.log")),
			zapcore.InfoLevel,
		),
		zap.AddCaller(),
	).Sugar()
	RequestLogger = zap.New(
		zapcore.NewCore(
			getEncoder(),
			getLogWriter(filepath.Join(path, "request.log")),
			zapcore.InfoLevel,
		),
		zap.AddCaller(),
	).Sugar()
	RepoLogger = zap.New(
		zapcore.NewCore(
			getEncoder(),
			getLogWriter(filepath.Join(path, "repo.log")),
			zapcore.InfoLevel,
		),
		zap.AddCaller(),
	).Sugar()
}
