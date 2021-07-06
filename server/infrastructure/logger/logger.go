package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type CustomLogger struct {
	zap.SugaredLogger
}

func Colored(color, text string) string {
	var iC uint8
	switch color {
	case "black":
		iC = 30
	case "red":
		iC = 31
	case "green":
		iC = 32
	case "yellow":
		iC = 33
	case "blue":
		iC = 34
	case "magenta":
		iC = 35
	case "cyan":
		iC = 36
	case "white":
		iC = 37
	case "r_black":
		iC = 40
	case "r_red":
		iC = 41
	case "r_green":
		iC = 42
	case "r_yellow":
		iC = 43
	case "r_blue":
		iC = 44
	case "r_magenta":
		iC = 45
	case "r_cyan":
		iC = 46
	case "r_white":
		iC = 47
	case "l_black":
		iC = 90
	case "l_red":
		iC = 91
	case "l_green":
		iC = 92
	case "l_yellow":
		iC = 93
	case "l_blue":
		iC = 94
	case "l_magenta":
		iC = 95
	case "l_cyan":
		iC = 96
	case "l_white":
		iC = 97
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", iC, text)
}

// NewSugarLogger ...
// DebugLevel  = -1
// InfoLevel   =  0
// WarnLevel   =  1
// ErrorLevel  =  2
// DPanicLevel =  3
// PanicLevel  =  4
// FatalLevel  =  5
func NewLogger(loglevel int, pretty bool) *CustomLogger {
	var (
		levelEnc  zapcore.LevelEncoder
		timeEnc   zapcore.TimeEncoder
		callerEnc zapcore.CallerEncoder
	)
	switch pretty {
	case true:
		levelEnc = zapcore.CapitalColorLevelEncoder
		timeEnc = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(Colored("l_black", t.Format("2006.01.02  15:04:05 .000")))
		}
		callerEnc = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(Colored("l_black", caller.TrimmedPath()))
		}
	case false:
		levelEnc = zapcore.CapitalLevelEncoder
		timeEnc = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006.01.02  15:04:05 .000"))
		}
		callerEnc = zapcore.ShortCallerEncoder
	}

	logger, err := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.Level(loglevel)),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  levelEnc,
			TimeKey:      "timestamp",
			EncodeTime:   timeEnc,
			CallerKey:    "caller",
			EncodeCaller: callerEnc,
		},
	}.Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	return &CustomLogger{
		SugaredLogger: *logger.Sugar(),
	}
}

func (cl *CustomLogger) FailOnError(err error, msg string) bool {
	if err != nil {
		cl.Errorf("%s: %s", msg, err)

		return true
	}

	return false
}
