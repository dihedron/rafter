package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

const (
	LogLevelOff = 100
)

var (
	configuration zap.Config
	restore       func()
)

func init() {
	var (
		err error
		exe string
	)

	// if the environment has a variable called after the application's
	// executable with the _DEBUG suffix (e.g. my-app -> MY_APP_DEBUG)
	// and it contains a valid value, set the log level accordingly;
	// otherwise disable the logging altogether.
	// TODO: it may be interesting to add a test on a special file to
	// configure/enable logging.
	var debugLevel zapcore.Level = LogLevelOff
	if exe, err = os.Executable(); err != nil {
		panic(fmt.Sprintf("error retrieving current application: %v", err))
	}
	debugVarName := strings.ReplaceAll(fmt.Sprintf("%s_DEBUG", strings.ToUpper(path.Base(exe))), "-", "_")
	if debugVarValue, present := os.LookupEnv(debugVarName); present {
		switch strings.ToLower(debugVarValue) {
		case "d", "dbg", "debug":
			debugLevel = zapcore.DebugLevel
		case "i", "inf", "info", "informational":
			debugLevel = zapcore.InfoLevel
		case "w", "wrn", "warn", "warning":
			debugLevel = zapcore.WarnLevel
		case "e", "err", "error":
			debugLevel = zapcore.ErrorLevel
		case "f", "ftl", "fatal":
			debugLevel = zapcore.PanicLevel
		case "off", "none":
			debugLevel = LogLevelOff
		}
	}

	if !strings.Contains(os.Args[0], ".test") && debugLevel != LogLevelOff {
		app := strings.Replace(filepath.Base(os.Args[0]), ".exe", "", 1)
		var (
			err    error
			logger *zap.Logger
		)
		if content, err := ioutil.ReadFile(app + ".json"); err == nil {
			if err := json.Unmarshal(content, &configuration); err == nil {

				// update the field tags to make elastic happy
				fillForElastic(&configuration)

				logger, err = configuration.Build()
				if err != nil {
					panic(fmt.Sprintf("error initialising logger: %v", err))
				}
				restore = zap.ReplaceGlobals(logger)
				zap.S().Info("application starting with custom log configuration")
				return
			}
		}

		configuration = zap.NewProductionConfig()
		configuration.Encoding = "json" // or "console"

		// update the field tags to make elastic happy
		fillForElastic(&configuration)

		configuration.OutputPaths = []string{fmt.Sprintf("%s-%d.log", app, os.Getpid())}
		configuration.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = configuration.Build()
		if err != nil {
			panic(fmt.Sprintf("error initialising logger: %v", err))
		}
		restore = zap.ReplaceGlobals(logger)

		SetLevel(debugLevel)

		zap.S().Info("application starting with default log configuration")
	}
}

// SetLevel sets the level of the logger globally.
func SetLevel(level zapcore.Level) {
	configuration.Level.SetLevel(level)
}

// EnableTestLogger is used to replace the current logger with
// one that writes to the one from the testing package when
// running Go tests.
func EnableTestLogger(t *testing.T) func() {
	//return zap.ReplaceGlobals(zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller(), zap.AddCallerSkip(1))))
	return zap.ReplaceGlobals(zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller())))
}

// RestoreGlobals restores the default logger ans sugared logger for Uber's zap.
func RestoreGlobals() {
	restore()
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
