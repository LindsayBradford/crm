// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

type AnnealerBuilder struct {
	annealer Annealer
}

func (this *AnnealerBuilder) SingleObjectiveAnnealer() *AnnealerBuilder {
	this.annealer = &singleObjectiveAnnealer{}
	this.annealer.Initialise()
	return this
}

func (this *AnnealerBuilder) WithStartingTemperature(temperature float64) *AnnealerBuilder {
	annealer := this.annealer
	annealer.setTemperature(temperature)
	return this
}

func (this *AnnealerBuilder) WithCoolingFactor(coolingFactor float64) *AnnealerBuilder {
	annealer := this.annealer
	annealer.setCoolingFactor(coolingFactor)
	return this
}

func (this *AnnealerBuilder) WithMaxIterations(iterations uint) *AnnealerBuilder {
	annealer := this.annealer
	annealer.setMaxIterations(iterations)
	return this
}

func (this *AnnealerBuilder) WithObservers(observers ...AnnealingObserver) *AnnealerBuilder {
	// PRE: observers supplied are not nil
	annealer := this.annealer

	for _, currObserver := range observers {
		annealer.AddObserver(currObserver)
	}

	return this
}

func (this *AnnealerBuilder) Build() Annealer {
	return this.annealer
}