package base

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/dihedron/rafter/logging"
)

type Base struct {
	Debug      string `short:"D" long:"debug" description:"The debug level of the application." optional:"yes" choice:"off" choice:"debug" choice:"info" choice:"warn" choice:"error"`
	CPUProfile string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes"`
	MemProfile string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes"`
	Logger     string `short:"L" long:"logger" description:"The logger to use." optional:"yes" choice:"zap" choice:"console" choice:"hcl" choice:"file" choice:"warn" choice:"off"`
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
