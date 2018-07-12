// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package profiling

import (
	"errors"
	"os"
	"runtime/pprof"

)

type NoParameterFunction func() error

// CpuProfileOfFunctionToFile establishes profiling based on what is passed as the cpuProfile
// parameter. If an empty string, it's assumed profiling is not needed.  A non-empty
// string is assumed to contain a path to a file in which profiling data is to be collated.

func CpuProfileOfFunctionToFile(functionToProfile NoParameterFunction, cpuProfilePath string) error {
	if cpuProfilePath == "" {
		return errors.New("empty cpuProfilePath supplied")
	}

	f, err := os.Create(cpuProfilePath)
	if err != nil {
		return err
	}

	pprof.StartCPUProfile(f)

	err = functionToProfile()

	defer pprof.StopCPUProfile()

	return err
}