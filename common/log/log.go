package log

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/common/telemetry"
	"go.uber.org/zap"
	"runtime"
	"strings"
)

type log struct {
	logProvider telemetry.Logger
}

const (
	skip = 2
)

var logInstance *log

func NewLogger(telemetryAPI *telemetry.API) {
	if logInstance == nil {
		logInstance = &log{
			logProvider: telemetryAPI.Logger(),
		}
	}
}

func Error(ctx context.Context, msg string, err error, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	fields = append(fields, fileField)
	logInstance.logProvider.Error(ctx, msg, err, convertFields(fields)...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	fields = append(fields, fileField)
	logInstance.logProvider.Info(ctx, msg, convertFields(fields)...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	fields = append(fields, fileField)
	logInstance.logProvider.Warn(ctx, msg, convertFields(fields)...)
}

func getFileAndLine() (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		i := strings.Index(file, "td")
		file = file[i:]
	}

	return file, line
}

func SetFieldString(param map[string]string) (fields []Field) {
	for i, val := range param {
		fields = append(fields, Field{Field: zap.String(i, val)})
	}
	return
}
