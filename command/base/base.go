package base

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/go-hclog"
)

type Base struct {
	Debug      string `short:"D" long:"debug" description:"The debug level of the application." optional:"yes" choice:"off" choice:"debug" choice:"info" choice:"warn" choice:"error"`
	CPUProfile string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes"`
	MemProfile string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes"`
	Logger     string `short:"L" long:"logger" description:"The logger to use." optional:"yes" default:"none" choice:"zap" choice:"console" choice:"hcl" choice:"file" choice:"log" choice:"none"`
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

func (cmd *Base) GetLogger(wrapped interface{}) logging.Logger {
	switch cmd.Logger {
	case "none":
		return &logging.NoOpLogger{}
	case "console":
		return logging.NewConsoleLogger(logging.StdOut)
	case "hcl":
		return logging.NewHCLLogger(wrapped.(hclog.Logger))
	case "zap":
		return logging.NewZapLogger()
	case "file":
		exe, _ := os.Executable()
		log := fmt.Sprintf("%s-%d.log", strings.Replace(exe, ".exe", "", -1), os.Getpid())
		return logging.NewFileLogger(log)
	}
	return &logging.NoOpLogger{}
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
