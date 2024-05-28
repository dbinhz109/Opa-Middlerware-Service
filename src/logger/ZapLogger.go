package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger
var commonFields = make([]zap.Field, 0, 16)

func init() {
	var err error
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	enccoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	enccoderConfig.StacktraceKey = "stackTrace"
	enccoderConfig.TimeKey = "timestamp"
	config.EncoderConfig = enccoderConfig

	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func InitLoggingData() {
	commonFields = append(
		commonFields,
		zap.String("service_", viper.GetString("application.name")),
		zap.String("ip_", viper.GetString("server.ip")),
		zap.String("port_", viper.GetString("server.port")),
	)
}

func Info(message string, fields ...zap.Field) {
	fields = append(fields, commonFields...)
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	fields = append(fields, commonFields...)
	zapLog.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	fields = append(fields, commonFields...)
	zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	fields = append(fields, commonFields...)
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	fields = append(fields, commonFields...)
	zapLog.Fatal(message, fields...)
}
