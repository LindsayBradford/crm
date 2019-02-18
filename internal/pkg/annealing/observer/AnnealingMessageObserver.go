// (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// AnnealingMessageObserver produces a stream of human-friendly, free-form text log entries from any observed
// Event instances received.
type AnnealingMessageObserver struct {
	AnnealingObserver
}

func (amo *AnnealingMessageObserver) WithLogHandler(handler logging.Logger) *AnnealingMessageObserver {
	amo.logHandler = handler
	return amo
}

func (amo *AnnealingMessageObserver) WithFilter(filter filters.Filter) *AnnealingMessageObserver {
	amo.filter = filter
	return amo
}

// ObserveEvent captures and converts Event instances into free-form text strings that it
// then passes onto its relevant Logger as an Info call.
func (amo *AnnealingMessageObserver) ObserveEvent(event observer.Event) {
	if amo.logHandler.BeingDiscarded(AnnealerLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	if observableAnnealer, isAnnealer := event.EventSource.(annealing.Observable); isAnnealer {
		annealer := wrapAnnealer(observableAnnealer)
		explorer := wrapSolutionExplorer(observableAnnealer.ObservableExplorer())

		var builder strings.FluentBuilder
		builder.Add("Id [", event.EventSource.Id(), "], ", "Event [", event.EventType.String(), "]: ")

		switch event.EventType {
		case observer.StartedAnnealing:
			builder.
				Add("Maximum Iterations [", annealer.MaximumIterations(), "], ").
				Add("Objective value [", explorer.ObjectiveValue(), "], ").
				Add("Temperature [", annealer.Temperature(), "], ").
				Add("Cooling Factor [", annealer.CoolingFactor(), "]")
		case observer.StartedIteration:
			builder.
				Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
				Add("Temperature [", annealer.Temperature(), "], ").
				Add("Objective value [", explorer.ObjectiveValue(), "]")
		case observer.FinishedIteration:
			builder.
				Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
				Add("Objective value [", explorer.ObjectiveValue(), "], ").
				Add("Change [", explorer.ChangeInObjectiveValue(), "], ").
				Add("Desirable? [", explorer.ChangeIsDesirable(), "], ").
				Add("Acceptance Probability [", explorer.AcceptanceProbability(), "], ").
				Add("Accepted? [", explorer.ChangeAccepted(), "]")
		case observer.FinishedAnnealing:
			builder.
				Add("Iteration [", annealer.CurrentIteration(), "/", annealer.MaximumIterations(), "], ").
				Add("Objective value [", explorer.ObjectiveValue(), "], ").
				Add("Temperature [", annealer.Temperature(), "]")
		case observer.Note:
			builder.Add("[", event.Note, "]")
		default:
			// deliberately does nothing extra
		}

		amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
	}
}
