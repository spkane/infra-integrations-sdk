package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SetupLogging redirects logs to stderr and configures the log level.
func SetupLogging(verbose bool) {
	logrus.SetOutput(os.Stderr)
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// New creates an already configured logger.
func New(verbose bool) *logrus.Logger {
	l := logrus.New()
	ConfigureLogger(l, verbose)

	return l
}

// ConfigureLogger configures an already created logger. Redirects logs to stderr and configures the log level.
func ConfigureLogger(logger *logrus.Logger, verbose bool) {
	logger.Out = os.Stderr
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

// Debug logs a formatted message at level Debug.
func Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info logs a formatted message at level Info.
func Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn logs a formatted message at level Warn.
func Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error logs a formatted message at level Error.
func Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal logs an error at level Fatal, and makes the program exit with an error code.
func Fatal(err error) {
	logrus.WithError(err).Fatal("can't continue")
}
