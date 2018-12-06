// Copyright (c) 2018 Australian Rivers Institute.

package model

import "github.com/LindsayBradford/crem/pkg/name"

type Model interface {
	name.Nameable
	DecisionVariableContainer

	TryRandomChange()
	AcceptChange()
	RevertChange()

	Clone() Model
}

// Container defines an interface embedding a Model
type Container interface {
	Model() Model
	SetModel(model Model)
}

type Contained struct {
	model Model
}

func (c *Contained) Model() Model {
	return c.model
}

func (c *Contained) SetModel(model Model) {
	c.model = model
}

type DecisionVariableContainer interface {
	DecisionVariable(name string) (DecisionVariable, error)
	DecisionVariableChange(decisionVariableName string) (float64, error)
}

var NullModel = new(nullModel)

func NewNullModel() *nullModel {
	newModel := new(nullModel).WithName("NullModel")
	return newModel
}

type nullModel struct {
	name.Named
}

func (nm *nullModel) WithName(name string) *nullModel {
	nm.SetName(name)
	return nm
}

func (nm *nullModel) TryRandomChange() {}
func (nm *nullModel) AcceptChange()    {}
func (nm *nullModel) RevertChange()    {}
func (nm *nullModel) DecisionVariable(name string) (DecisionVariable, error) {
	newVariable := DecisionVariableImpl{
		name:  name,
		value: 0,
	}
	return &newVariable, nil
}
func (nm *nullModel) DecisionVariableChange(decisionVariableName string) (float64, error) {
	return 0, nil
}
func (nm *nullModel) SetDecisionVariable(name string, value float64) error { return nil }
func (nm *nullModel) Clone() Model {
	clone := *nm
	return &clone
}