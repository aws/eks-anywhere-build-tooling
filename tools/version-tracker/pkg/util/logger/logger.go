package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	pkgLogger logr.Logger = logr.Discard()
	Verbosity int
)

// Init initializes the package logger. Repeat calls will overwrite the package logger which may
// result in unexpected behavior.
func Init(verbosityLevel int) error {
	Verbosity = verbosityLevel
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = nil
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}

	// Level 6 and above are used for debugging and we want a different log structure for debug
	// logs.
	if verbosityLevel >= 6 {
		encoderCfg.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			// Because we use negated levels it is necessary to negate the level again so the
			// output appears in a V0 format.
			//
			// See logrAtomicLevel().
			enc.AppendString(fmt.Sprintf("V%d", -int(l)))
		}
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Build the encoder and logger.
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logrAtomicLevel(verbosityLevel)),
	)
	logger := zap.New(core)

	// Configure package state so the logger can be used by other packages.

	pkgLogger = zapr.NewLogger(logger)

	return nil
}

// logrAtomicLevel creates a zapcore.AtomicLevel compatible with go-logr.
func logrAtomicLevel(level int) zap.AtomicLevel {
	// The go-logr wrapper uses custom Zap log levels. To represent this in Zap, its
	// necessary to negate the level to circumvent Zap level constraints.
	//
	// See https://github.com/go-logr/zapr/blob/master/zapr.go#L50.
	return zap.NewAtomicLevelAt(zapcore.Level(-level))
}

// Get returns the logger instance that has been previously set.
// If no logger has been set, it returns a null logger.
func Get() logr.Logger {
	return pkgLogger
}

// Info logs a non-error message with the given key/value pairs as context.
//
// The msg argument should be used to add some constant description to
// the log line. The key/value pairs can then be used to add additional
// variable information. The key/value pairs should alternate string
// keys and arbitrary values.
func Info(msg string, keysAndValues ...interface{}) {
	Get().Info(msg, keysAndValues...)
}

// V returns an Logger value for a specific verbosity level, relative to
// this Logger. In other words, V values are additive.  V higher verbosity
// level means a log message is less important.  It's illegal to pass a log
// level less than zero.
func V(level int) logr.Logger {
	return Get().V(level)
}
