// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import (
	"errors"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	rand2 "github.com/LindsayBradford/crem/pkg/rand"
)

type Model struct {
	name                  string
	randomNumberGenerator *rand2.ConcurrencySafeRand

	decisionVariables     map[string]model.DecisionVariable
	tempDecisionVariables map[string]model.DecisionVariable
	parameters            Parameters
}

func New() *Model {
	newModel := new(Model)
	newModel.name = "DumbModel"

	newModel.randomNumberGenerator = rand2.NewTimeSeeded()
	newModel.buildDecisionVariables()
	newModel.parameters.Initialise()

	initialValue := newModel.parameters.GetFloat64(InitialObjectiveValue)
	newModel.decisionVariables[model.ObjectiveValue].SetValue(initialValue)
	newModel.tempDecisionVariables[model.ObjectiveValue].SetValue(initialValue)

	return newModel
}

func (dm *Model) Name() string { return dm.name }
func (dm *Model) WithName(name string) *Model {
	dm.name = name
	return dm
}
func (dm *Model) SetName(name string) { dm.name = name }

func (dm *Model) buildDecisionVariables() {
	dm.decisionVariables = make(map[string]model.DecisionVariable, 1)
	dm.tempDecisionVariables = make(map[string]model.DecisionVariable, 1)

	objectiveValueVar := new(model.DecisionVariableImpl)
	objectiveValueVar.SetName(model.ObjectiveValue)
	dm.decisionVariables[model.ObjectiveValue] = objectiveValueVar

	tempObjectiveValueVar := new(model.DecisionVariableImpl)
	tempObjectiveValueVar.SetName(model.ObjectiveValue)
	dm.tempDecisionVariables[model.ObjectiveValue] = tempObjectiveValueVar
}

func (dm *Model) WithParameters(params parameters.Map) *Model {
	dm.parameters.Merge(params)

	initialValue := dm.parameters.GetFloat64(InitialObjectiveValue)
	dm.decisionVariables[model.ObjectiveValue].SetValue(initialValue)
	dm.tempDecisionVariables[model.ObjectiveValue].SetValue(initialValue)

	return dm
}

func (dm *Model) ParameterErrors() error {
	return dm.parameters.ValidationErrors()
}

const (
	downward = -1
	upward   = 1
)

func (dm *Model) TryRandomChange() {
	originalValue := dm.objectiveValue()
	change := dm.generateRandomChange()
	newValue := dm.capChangeOverRange(originalValue + change)
	dm.setObjectiveValue(newValue)
}

func (dm *Model) generateRandomChange() float64 {
	randomValue := dm.randomNumberGenerator.Intn(2)

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
	return dm.decisionVariables[model.ObjectiveValue].Value()
}

func (dm *Model) setObjectiveValue(value float64) {
	dm.tempDecisionVariables[model.ObjectiveValue].SetValue(value)
}

func (dm *Model) copyDecisionVarValueToTemp(varName string) {
	dm.tempDecisionVariables[varName].SetValue(
		dm.decisionVariables[varName].Value(),
	)
}

func (dm *Model) copyTempDecisionVarValueToActual(varName string) {
	dm.decisionVariables[varName].SetValue(
		dm.tempDecisionVariables[varName].Value(),
	)
}

func (dm *Model) AcceptChange() {
	dm.copyTempDecisionVarValueToActual(model.ObjectiveValue)
}

func (dm *Model) RevertChange() {
	dm.copyDecisionVarValueToTemp(model.ObjectiveValue)
}

func (dm *Model) DecisionVariable(name string) (model.DecisionVariable, error) {
	if variable, found := dm.decisionVariables[name]; found == true {
		return variable, nil
	}
	return model.NullDecisionVariable, errors.New("decision variable [" + name + "] not defined for model [" + dm.Name() + " ].")
}

func (dm *Model) DecisionVariableChange(variableName string) (float64, error) {
	decisionVariable, foundActual := dm.decisionVariables[variableName]
	tmpDecisionVar, foundTemp := dm.tempDecisionVariables[decisionVariable.Name()]
	if !foundActual || !foundTemp {
		return 0, errors.New("no temporary decision variable of name [" + decisionVariable.Name() + "] in model [" + dm.Name() + "].")
	}

	difference := tmpDecisionVar.Value() - decisionVariable.Value()
	return difference, nil
}

func (dm *Model) Clone() model.Model {
	clone := *dm
	return &clone
}

const (
	InitialObjectiveValue = "InitialObjectiveValue"
	MinimumObjectiveValue = "MinimumObjectiveValue"
	MaximumObjectiveValue = "MaximumObjectiveValue"
)

type Parameters struct {
	parameters.Parameters
}

func (kp *Parameters) Initialise() *Parameters {
	kp.Parameters.Initialise()
	kp.buildMetaData()
	kp.CreateDefaults()
	return kp
}

func (kp *Parameters) buildMetaData() {
	kp.AddMetaData(
		parameters.MetaData{
			Key:          InitialObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(1000),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          MinimumObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(0),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          MaximumObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(2000),
		},
	)
}
