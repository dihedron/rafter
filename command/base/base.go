package base

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/dihedron/rafter/logging"
	"github.com/dihedron/rafter/logging/console"
	"github.com/dihedron/rafter/logging/file"
	"github.com/dihedron/rafter/logging/noop"
	"github.com/dihedron/rafter/logging/uber"
)

type Base struct {
	Debug      string `short:"D" long:"debug" description:"The debug level of the application." optional:"yes" choice:"off" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" default:"debug"`
	CPUProfile string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes"`
	MemProfile string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes"`
	Logger     string `short:"L" long:"logger" description:"The logger to use." optional:"yes" default:"none" choice:"zap" choice:"console" choice:"file" choice:"log" choice:"none"`
}

type Closer struct {
	file *os.File
}

func (c *Closer) Close() {
	if c.file != nil {
		pprof.StopCPUProfile()
		c.file.Close()
	}
}

func (cmd *Base) GetLogger() logging.Logger {
	switch cmd.Debug {
	case "trace":
		logging.SetLevel(logging.LevelTrace)
	case "debug":
		logging.SetLevel(logging.LevelDebug)
	case "info":
		logging.SetLevel(logging.LevelInfo)
	case "warn":
		logging.SetLevel(logging.LevelWarn)
	case "error":
		logging.SetLevel(logging.LevelError)
	case "off":
		logging.SetLevel(logging.LevelOff)
	}
	switch cmd.Logger {
	case "none":
		return &noop.Logger{}
	case "console":
		return console.NewLogger(console.StdOut)
	case "zap":
		logger, _ := uber.NewLogger()
		return logger
	case "file":
		exe, _ := os.Executable()
		log := fmt.Sprintf("%s-%d.log", strings.Replace(exe, ".exe", "", -1), os.Getpid())
		return file.NewLogger(log)
	}
	return &noop.Logger{}
}

func (cmd *Base) ProfileCPU(logger logging.Logger) *Closer {
	var f *os.File
	if cmd.CPUProfile != "" {
		var err error
		f, err = os.Create(cmd.CPUProfile)
		if err != nil {
			logger.Error("could not create CPU profile at %s: %v", cmd.CPUProfile, err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Error("could not start CPU profiler: %v", err)
		}
	}
	return &Closer{
		file: f,
	}
}

func (cmd *Base) ProfileMemory(logger logging.Logger) {
	if cmd.MemProfile != "" {
		f, err := os.Create(cmd.MemProfile)
		if err != nil {
			logger.Error("could not create memory profile at %s: %v", cmd.MemProfile, err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			logger.Error("could not write memory profile: %v", err)
		}
	}
}
