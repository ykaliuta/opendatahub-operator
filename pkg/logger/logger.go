package logger

import (
	"os"

	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func NewLogger(mode string) logr.Logger {
	raw := NewRawLogger(mode).GetSink()
	sink := NewSink(raw)
	return logr.New(sink)
}

func UpdateLogger(l logr.Logger, mode string) {
	logrSink := l.GetSink()
	sink, ok := logrSink.(*Sink)
	if !ok {
		l.Info("logger sink is not of type *Sink")
		return
	}

	raw := NewRawLogger(mode).GetSink()
	sink.SetSink(raw)
}

// in DSC component, to use different mode for logging, e.g. development, production
// when not set mode it falls to "default" which is used by startup main.go.
func NewRawLogger(mode string) logr.Logger {
	var opts zap.Options
	switch mode {
	case "devel", "development": //  the most logging verbosity
		opts = zap.Options{
			Development:     true,
			StacktraceLevel: zapcore.WarnLevel,
			Level:           zapcore.DebugLevel,
			DestWriter:      os.Stdout,
		}
	case "prod", "production": // the least logging verbosity
		opts = zap.Options{
			Development:     false,
			StacktraceLevel: zapcore.ErrorLevel,
			Level:           zapcore.ErrorLevel,
			DestWriter:      os.Stdout,
			EncoderConfigOptions: []zap.EncoderConfigOption{func(config *zapcore.EncoderConfig) {
				config.EncodeTime = zapcore.ISO8601TimeEncoder // human readable not epoch
				config.EncodeDuration = zapcore.SecondsDurationEncoder
				config.LevelKey = "LogLevel"
				config.NameKey = "Log"
				config.CallerKey = "Caller"
				config.MessageKey = "Message"
				config.TimeKey = "Time"
				config.StacktraceKey = "Stacktrace"
			}},
		}
	default:
		opts = zap.Options{
			Development:     false,
			StacktraceLevel: zapcore.ErrorLevel,
			Level:           zapcore.InfoLevel,
			DestWriter:      os.Stdout,
		}
	}
	return zap.New(zap.UseFlagOptions(&opts))
}
