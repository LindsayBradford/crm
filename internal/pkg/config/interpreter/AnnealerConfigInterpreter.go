// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/cooling/coolants/averaged"
	coolingSuppapitnarm "github.com/LindsayBradford/crem/internal/pkg/annealing/cooling/coolants/suppapitnarm"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/kirkpatrick"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/suppapitnarm"
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/pkg/errors"
)

type AnnealerConfigInterpreter struct {
	errors              *compositeErrors.CompositeError
	registeredAnnealers map[data.AnnealerType]AnnealerConfigFunction

	annealer annealing.Annealer
}

func (i *AnnealerConfigInterpreter) initialise() *AnnealerConfigInterpreter {
	i.registeredAnnealers = make(map[data.AnnealerType]AnnealerConfigFunction, 0)
	i.errors = compositeErrors.New("Annealer Configuration")
	i.annealer = &annealers.NullAnnealer{}
	return i
}

type AnnealerConfigFunction func(config data.AnnealerConfig) annealing.Annealer

func NewAnnealerConfigInterpreter() *AnnealerConfigInterpreter {
	newInterpreter := new(AnnealerConfigInterpreter).initialise().
		RegisteringAnnealer(
			data.UnspecifiedAnnealerType,
			func(config data.AnnealerConfig) annealing.Annealer {
				newAnnealer := new(annealers.NullAnnealer)
				newAnnealer.Initialise()
				return newAnnealer
			},
		).RegisteringAnnealer(
		data.Kirkpatrick,
		func(config data.AnnealerConfig) annealing.Annealer {
			newAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)
			newAnnealer.Initialise()

			newExplorer := kirkpatrick.New()
			newAnnealer.SetSolutionExplorer(newExplorer)

			newAnnealer.SetParameters(config.Parameters)

			return newAnnealer
		},
	).RegisteringAnnealer(
		data.Suppapitnarm,
		func(config data.AnnealerConfig) annealing.Annealer {
			newAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)
			newAnnealer.Initialise()

			newExplorer := suppapitnarm.New().WithCoolant(coolingSuppapitnarm.NewCoolant())

			newAnnealer.SetSolutionExplorer(newExplorer)

			newAnnealer.SetParameters(config.Parameters)

			return newAnnealer
		},
	).RegisteringAnnealer(
		data.AveragedSuppapitnarm,
		func(config data.AnnealerConfig) annealing.Annealer {
			newAnnealer := new(annealers.ElapsedTimeTrackingAnnealer)
			newAnnealer.Initialise()

			newExplorer := suppapitnarm.New().WithCoolant(averaged.NewCoolant())

			newAnnealer.SetSolutionExplorer(newExplorer)

			newAnnealer.SetParameters(config.Parameters)

			return newAnnealer
		},
	)
	return newInterpreter
}

func (i *AnnealerConfigInterpreter) Interpret(config *data.AnnealerConfig) *AnnealerConfigInterpreter {
	assert.That(config != nil)

	if _, foundAnnealer := i.registeredAnnealers[config.Type]; !foundAnnealer {
		i.errors.Add(
			errors.New("configuration specifies annealer type [\"" +
				config.Type.String() + "\"], but no annealers are registered for that type"),
		)
		return i
	}

	configFunction := i.registeredAnnealers[config.Type]
	newAnnealer := configFunction(*config)
	if parameterisedModel, hasParameters := newAnnealer.(parameters.Container); hasParameters {
		if paramErrors := parameterisedModel.ParameterErrors(); paramErrors != nil {
			wrappedErrors := errors.Wrap(paramErrors, "building annealer ["+config.Type.String()+"]")
			i.errors.Add(wrappedErrors)
			return i
		}
	}
	i.annealer = newAnnealer
	return i
}

func (i *AnnealerConfigInterpreter) RegisteringAnnealer(annealerType data.AnnealerType, configFunction AnnealerConfigFunction) *AnnealerConfigInterpreter {
	i.registeredAnnealers[annealerType] = configFunction
	return i
}

func (i *AnnealerConfigInterpreter) Annealer() annealing.Annealer {
	return i.annealer
}

func (i *AnnealerConfigInterpreter) Errors() error {
	if i.errors.Size() > 0 {
		return i.errors
	}
	return nil
}
