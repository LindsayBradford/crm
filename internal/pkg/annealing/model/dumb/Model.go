// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"math"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/variables"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Model struct {
	name.NameContainer
	name.IdentifiableContainer
	rand.RandContainer

	parameters Parameters
	variables.ContainedDecisionVariables
}

func New() *Model {
	newModel := new(Model)
	newModel.SetName("DumbModel")

	newModel.DecisionVariables()
	newModel.parameters.Initialise()

	newModel.ContainedDecisionVariables.Initialise()
	newModel.ContainedDecisionVariables.NewForName(variable.ObjectiveValue)

	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.ContainedDecisionVariables.SetValue(variable.ObjectiveValue, initialValue)

	return newModel
}

func (dm *Model) WithName(name string) *Model {
	dm.SetName(name)
	return dm
}

func (dm *Model) WithParameters(params parameters.Map) *Model {
	dm.parameters.Merge(params)

	initialValue := dm.parameters.GetFloat64(InitialObjectiveValue)
	dm.ContainedDecisionVariables.SetValue(variable.ObjectiveValue, initialValue)

	return dm
}

func (dm *Model) ParameterErrors() error {
	return dm.parameters.ValidationErrors()
}

const (
	downward = -1
	upward   = 1
)

func (dm *Model) Initialise() {
	dm.SetRandomNumberGenerator(rand.NewTimeSeeded())
}

func (dm *Model) TearDown() {
	// This model doesn't need any special tearDown.
}

func (dm *Model) TryRandomChange() {
	originalValue := dm.objectiveValue()
	change := dm.generateRandomChange()
	newValue := dm.capChangeOverRange(originalValue + change)
	dm.setObjectiveValue(newValue)
}

func (dm *Model) generateRandomChange() float64 {
	randomValue := dm.RandomNumberGenerator().Intn(2)

	var changeInObjectiveValue float64
	switch randomValue {
	case 0:
		changeInObjectiveValue = downward
	case 1:
		changeInObjectiveValue = upward
	}

	return changeInObjectiveValue
}

func (dm *Model) capChangeOverRange(value float64) float64 {
	maxCappedValue := math.Max(dm.parameters.GetFloat64(MinimumObjectiveValue), value)
	bothCappedValue := math.Min(dm.parameters.GetFloat64(MaximumObjectiveValue), maxCappedValue)
	return bothCappedValue
}

func (dm *Model) objectiveValue() float64 {
	return dm.ContainedDecisionVariables.Value(variable.ObjectiveValue)
}

func (dm *Model) setObjectiveValue(value float64) {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).SetInductiveValue(value)
}

func (dm *Model) SetDecisionVariable(name string, value float64) {
	dm.ContainedDecisionVariables.SetValue(name, value)
}

func (dm *Model) AcceptChange() {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).AcceptInductiveValue()
}

func (dm *Model) RevertChange() {
	dm.ContainedDecisionVariables.Variable(variable.ObjectiveValue).RejectInductiveValue()
}

func (dm *Model) ActiveManagementActions() []action.ManagementAction { return nil }
func (dm *Model) PlanningUnits() solution.PlanningUnitIds            { return nil }

func (dm *Model) DeepClone() model.Model {
	clone := *dm
	clone.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
