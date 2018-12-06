// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"math/rand"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Explorer interface {
	Observable
	name.Nameable

	ScenarioId() string
	SetScenarioId(name string)

	Initialise()
	TearDown()

	Clone() Explorer
	TryRandomChange(temperature float64)

	SetRandomNumberGenerator(*rand.Rand)
	RandomNumberGenerator() *rand.Rand

	model.Container
	logging.Container
}

type Observable interface {
	ObjectiveValue() float64
	ChangeInObjectiveValue() float64
	AcceptanceProbability() float64
	ChangeIsDesirable() bool
	ChangeAccepted() bool
}

type Container interface {
	SolutionExplorer() Explorer
	SetSolutionExplorer(explorer Explorer) error
}
