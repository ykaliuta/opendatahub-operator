package logger

import (
	"os"

	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// NewNamedLogger creates a new logger for a component.
// If the mode is set (so can be different from the default one),
// it will create a new logger with the specified mode's options.
func NewNamedLogger(log logr.Logger, name string, mode string) logr.Logger {
	if mode != "" {
		log = NewLogger(mode)
	}
	return log.WithName(name)
}

// in DSC component, to use different mode for logging, e.g. development, production
// when not set mode it falls to "default" which is used by startup main.go.
func NewLogger(mode string) logr.Logger {
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
