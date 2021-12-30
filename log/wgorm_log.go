package log

import (
	"context"
	"fmt"
	"log"
)

var Logger = WgormLogger{}

type WgormLogger struct {
}

var (
	red    = string([]byte{27, 91, 57, 49, 109})
	yellow = string([]byte{27, 91, 57, 51, 109})
	blue   = string([]byte{27, 91, 57, 52, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

func (w *WgormLogger) InfoOf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf("[Info] "+format, args...)
	log.Println(blue, msg, reset)
}

func (w *WgormLogger) WarnOf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf("[Warn] "+format, args...)
	log.Println(yellow, msg, reset)
}

func (w *WgormLogger) ErrorOf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf("[Error] "+format, args...)
	log.Println(red, msg, reset)
}
