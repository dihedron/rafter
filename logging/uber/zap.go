package uber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dihedron/rafter/logging"
	"go.uber.org/zap"
)

// Logger is an adapter that allows to log using Uber's Zap
// wherever a Logger interface is expected.
type Logger struct {
	logger *zap.Logger
}

var (
	configuration zap.Config
	restore       func()
)

// NewLogger initialises a Zap logger, either by locating and loading
// a configuration fil from disk, or by assuming the sane defaults
// for a production environment.
func NewLogger() (*Logger, error) {

	// check if there's a file called brokerd-log.json aside the
	// application excutable; if so, load it as it contains the
	// logger configuration; if not, assume default for production
	app := strings.Replace(filepath.Base(os.Args[0]), ".exe", "", 1)
	content, err := ioutil.ReadFile(app + "-log.json")
	if err == nil { // the file exists
		if err := json.Unmarshal(content, &configuration); err != nil {
			return nil, fmt.Errorf("error unmarshalling log configuration from '%s': %w", app+"-log.json", err)
		}
		// update the field tags to make Elastic happy
		fillForElastic(&configuration)
		logger, err := configuration.Build()
		if err != nil {
			return nil, fmt.Errorf("error bulding logging configuration: %w", err)
		}
		restore = zap.ReplaceGlobals(logger)
		zap.S().Info("application starting with custom log configuration")
		return &Logger{
			logger: logger.WithOptions(zap.AddCallerSkip(1)),
			// logger: logger,
		}, nil
	}
	// configuration does not exist, use default
	configuration = zap.NewProductionConfig()
	configuration.Encoding = "json" // or "console"
	// update the field tags to make elastic happy
	fillForElastic(&configuration)
	configuration.OutputPaths = []string{fmt.Sprintf("%s-%d.log", app, os.Getpid())}
	configuration.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := configuration.Build()
	if err != nil {
		return nil, fmt.Errorf("error initialising logger: %w", err)
	}
	restore = zap.ReplaceGlobals(logger)
	zap.S().Info("application starting with default log configuration")

	return &Logger{
		logger: logger.WithOptions(zap.AddCallerSkip(1)),
		//logger: logger,
	}, nil
}

// Trace logs a message at LevelTrace level.
func (l *Logger) Trace(format string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelTrace {
		l.logger.Sugar().Debugf(format, args...)
	}
}

// Debug logs a message at LevelDebug level.
func (l *Logger) Debug(format string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelDebug {
		l.logger.Sugar().Debugf(format, args...)
	}
}

// Info logs a message at LevelInfo level.
func (l *Logger) Info(format string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelInfo {
		l.logger.Sugar().Infof(format, args...)
	}
}

// Warn logs a message at LevelWarn level.
func (l *Logger) Warn(format string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelWarn {
		l.logger.Sugar().Warnf(format, args...)
	}
}

// Error logs a message at LevelError level.
func (l *Logger) Error(format string, args ...interface{}) {
	if logging.GetLevel() <= logging.LevelError {
		l.logger.Sugar().Errorf(format, args...)
	}
}

func fillForElastic(configuration *zap.Config) {
	// configuration.EncoderConfig.MessageKey = "message"
	// configuration.EncoderConfig.LevelKey = "log.level"
	// configuration.EncoderConfig.TimeKey = "@timestamp"
	// configuration.EncoderConfig.NameKey = "log.logger"
	// configuration.EncoderConfig.CallerKey = "log.origin.file.name"
	// configuration.EncoderConfig.StacktraceKey = "error.stack_trace"
	// configuration.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	// configuration.InitialFields = map[string]interface{}{
	// 	"service.name":        appinfo.ServiceName,
	// 	"service.version":     fmt.Sprintf("v%s@%s", appinfo.GitTag, appinfo.GitCommit),
	// 	"service.environment": os.Getenv("BROKERD_STAGE"),
	// }
	// if configuration.InitialFields["service.environment"] == "" {
	// 	configuration.InitialFields["service.environment"] = "development"
	// }
}
