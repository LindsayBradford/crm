// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

type singleObjectiveAnnealer struct {
	annealerBase
}

func (this *singleObjectiveAnnealer) Anneal() {
	this.notifyObserversWith("I'm a single-objectve annealer")
	this.annealerBase.Anneal()
}