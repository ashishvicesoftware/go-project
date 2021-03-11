package logWriter

import (
	"os"
	"io"
	"time"
	"fmt"

	"github.com/jsternberg/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var production bool

type WriteSyncer struct {
    io.Writer
}

func (ws WriteSyncer) Sync() error {
    return nil
}

func SetProductionFlag(prod bool) {
	production = prod
}

func init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339))
	}
	LogToFile()
	var ioWriter = &lumberjack.Logger{
        Filename:   "gologfile",
        MaxSize:    10, // MB
        MaxBackups: 10,  // number of backups
        MaxAge:     28, //days
        LocalTime:  true,
        Compress:   false, // disabled by default
	}
	var sw = WriteSyncer{
        ioWriter,
    }
	
	// if err != nil {
	// 	logger = zap.New(zapcore.NewCore(
	// 		zaplogfmt.NewEncoder(config),
	// 		zapcore.AddSync(os.Stdout),
	// 		zapcore.DebugLevel))
	// 	logger.Error("Could not open log file gologfile")
	// } else {
		logger = zap.New(zapcore.NewTee(zapcore.NewCore(
				zaplogfmt.NewEncoder(config),
				zapcore.AddSync(os.Stdout),
				zapcore.DebugLevel),
				// zapcore.NewCore(
				// 	zaplogfmt.NewEncoder(config),
				// 	zapcore.AddSync(f),
				// 	zapcore.DebugLevel),
				zapcore.NewCore(
					zaplogfmt.NewEncoder(config), 
					sw, 
					zapcore.DebugLevel),
				))	
	// }

	// core := zapcore.NewTee(
    //     zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
    //     zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
    // )

    // l := zap.New(core)
}
func Debug(msg string, fields ...zap.Field) {
	if production == false {
		logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(err error) {
	logger.Fatal(err.Error())
}

func LogToFile() (io.Writer, error) {
	f, err := os.OpenFile("gologfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Error opening file: %v", err.Error())
	}
	defer f.Close()
	return f, err
}

// func getWriteSyncer(logName string) zapcore.WriteSyncer {
//     var ioWriter = &lumberjack.Logger{
//         Filename:   logName,
//         MaxSize:    10, // MB
//         MaxBackups: 3,  // number of backups
//         MaxAge:     28, //days
//         LocalTime:  true,
//         Compress:   false, // disabled by default
//     }
//     var sw = WriteSyncer{
//         ioWriter,
//     }
//     return sw
// }