package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type severity string

// ref: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry?hl=ja#LogSeverity
const (
	severityInfo     severity = "INFO"
	severityError    severity = "ERROR"
	severityWarn     severity = "WARNING"
	severityCritical severity = "CRITICAL"
	severityDebug    severity = "DEBUG"
)

func LogDebugf(ctx context.Context, format string, v ...interface{}) {
	debugf(ctx, fmt.Sprintf(format, v...))
}

func LogDebug(ctx context.Context, v interface{}) {
	debugf(ctx, fmt.Sprint(v))
}

func debugf(ctx context.Context, msg string) {
	logPrintf(ctx, severityDebug, msg)
}

func LogInfof(ctx context.Context, format string, v ...interface{}) {
	infof(ctx, fmt.Sprintf(format, v...))
}

func LogInfo(ctx context.Context, v interface{}) {
	infof(ctx, fmt.Sprint(v))
}

func infof(ctx context.Context, msg string) {
	logPrintf(ctx, severityInfo, msg)
}

func LogErrorf(ctx context.Context, format string, v ...interface{}) {
	errorf(ctx, fmt.Sprintf(format, v...))
}

func LogError(ctx context.Context, v interface{}) {
	errorf(ctx, fmt.Sprint(v))
}

func errorf(ctx context.Context, msg string) {
	logPrintf(ctx, severityError, msg)
}

func LogWarningf(ctx context.Context, format string, v ...interface{}) {
	warningf(ctx, fmt.Sprintf(format, v...))
}

func LogWarning(ctx context.Context, v interface{}) {
	warningf(ctx, fmt.Sprint(v))
}

func warningf(ctx context.Context, msg string) {
	logPrintf(ctx, severityWarn, msg)
}

func LogCriticalf(ctx context.Context, format string, v ...interface{}) {
	criticalf(ctx, fmt.Sprintf(format, v...))
}

func LogCritical(ctx context.Context, v interface{}) {
	criticalf(ctx, fmt.Sprint(v))
}

func criticalf(ctx context.Context, msg string) {
	logPrintf(ctx, severityCritical, msg)
}

func logPrintf(ctx context.Context, s severity, msg string) {
	logger := log.New(getWriter(), "", 0)

	// The info for struct logging are listed below.

	// ref: https://cloud.google.com/logging/docs/agent/configuration?hl=ja#special-fields
	entry := map[string]interface{}{
		"message":  msg,
		"severity": s,
	}
	payload, _ := json.Marshal(entry)
	logger.Println(string(payload))
}

var getWriter = func() io.Writer {
	return os.Stdout
}
