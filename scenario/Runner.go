// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"fmt"
	"sync"
	. "time"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/logging"
)

type CallableRunner interface {
	Run() error
	LogHandler() logging.Logger
}

type Runner struct {
	annealer   annealing.Annealer
	logHandler logging.Logger

	name              string
	operationType     string
	runNumber         uint64
	maxConcurrentRuns uint64
	tearDown          func()

	startTime  Time
	finishTime Time
}

func (runner *Runner) ForAnnealer(annealer annealing.Annealer) *Runner {
	runner.runNumber = 1
	runner.maxConcurrentRuns = 1 // Sequential by default
	runner.name = "Default Scenario"
	runner.tearDown = defaultTearDown

	runner.logHandler = annealer.LogHandler()
	runner.annealer = annealer

	return runner
}

func (runner *Runner) WithName(name string) *Runner {
	if name != "" {
		runner.name = name
	}
	return runner
}

func (runner *Runner) WithRunNumber(runNumber uint64) *Runner {
	if runNumber > 0 {
		runner.runNumber = runNumber
	}
	return runner
}

func (runner *Runner) WithTearDownFunction(tearDown func()) *Runner {
	if tearDown != nil {
		runner.tearDown = tearDown
	}
	return runner
}

func defaultTearDown() {
	// deliberately does nothing
}

func (runner *Runner) WithMaximumConcurrentRuns(maxConcurrentRuns uint64) *Runner {
	if maxConcurrentRuns > 0 {
		runner.maxConcurrentRuns = maxConcurrentRuns
	}
	return runner
}

func (runner *Runner) Run() error {
	runner.logScenarioStartMessage()
	runner.startTime = Now()

	runError := runner.runScenario()

	runner.finishTime = Now()
	runner.logHandler.Info("Finished running scenario \"" + runner.name + "\"")

	runner.logHandler.Info(runner.generateElapsedTimeString())

	runner.tearDown()

	return runError
}

func (runner *Runner) LogHandler() logging.Logger {
	return runner.logHandler
}

func (runner *Runner) logScenarioStartMessage() {
	var runTypeText string
	if runner.maxConcurrentRuns > 1 {
		runTypeText = fmt.Sprintf("executing a maximum of %d runs concurrently.", runner.maxConcurrentRuns)
	} else {
		runTypeText = "executing runs sequentially"
	}

	message := fmt.Sprintf("Scenario \"%s\": configured for %d run(s), %s", runner.name, runner.runNumber, runTypeText)
	runner.logHandler.Info(message)
}

func (runner *Runner) generateElapsedTimeString() string {
	return fmt.Sprintf("Total elapsed time of scenario = [%v]", runner.ElapsedTime())
}

func (runner *Runner) ElapsedTime() Duration {
	return runner.finishTime.Sub(runner.startTime)
}

func (runner *Runner) runScenario() error {
	var runWaitGroup sync.WaitGroup

	concurrentRunGuard := make(chan struct{}, runner.maxConcurrentRuns)

	doRun := func(runNumber uint64) {
		runner.run(runNumber)
		<-concurrentRunGuard
		runWaitGroup.Done()
	}

	runWaitGroup.Add(int(runner.runNumber))

	for runNumber := uint64(1); runNumber <= runner.runNumber; runNumber++ {
		concurrentRunGuard <- struct{}{}
		go doRun(runNumber)
	}

	runWaitGroup.Wait()

	return nil
}

func (runner *Runner) run(runNumber uint64) {
	annealerCopy := runner.annealer.Clone()

	annealerCopy.SetId(runner.generateCloneId(runNumber))

	runner.logRunStartMessage(runNumber)
	annealerCopy.Anneal()
	runner.logRunFinishedMessage(runNumber)
}

func (runner *Runner) generateCloneId(runNumber uint64) string {
	if runner.runNumber > 1 {
		return fmt.Sprintf("%s (%d/%d)", runner.name, runNumber, runner.runNumber)
	} else {
		return runner.name
	}
}

func (runner *Runner) logRunStartMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run started")
	}
}

func (runner *Runner) logRunFinishedMessage(runNumber uint64) {
	if runner.runNumber > 1 {
		runner.logHandler.Info(runner.generateCloneId(runNumber) + ": run finished")
	}
}
