// Copyright (c) 2018 Australian Rivers Institute.

package filters

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/observer"
)

// PercentileOfIterationsPerAnnealingFilter filters FinishedIteration Annealing Event instances at a rate of 1 every
// percentile number of iterations received. . StartedIteration events are completely filtered out. All other event types are allowed through to the LogHandler.
type PercentileOfIterationsPerAnnealingFilter struct {
	MaximumIterations  uint64
	percentileToReport float64
	iterationModulo    uint64
	generatesEvents    bool
}

// WithPercentileOfIterations defines the number of Annealing Iteration Event instances to report over the entire run.
func (m *PercentileOfIterationsPerAnnealingFilter) WithPercentileOfIterations(percentile float64) *PercentileOfIterationsPerAnnealingFilter {
	m.percentileToReport = percentile / 100
	m.deriveModuloIfPossible()
	return m
}

// WithPercentileOfIterations defines the number of Annealing Iteration Event instances to report over the entire run.
func (m *PercentileOfIterationsPerAnnealingFilter) WithMaximumIterations(MaximumIterations uint64) *PercentileOfIterationsPerAnnealingFilter {
	m.MaximumIterations = MaximumIterations
	m.deriveModuloIfPossible()
	return m
}

func (m *PercentileOfIterationsPerAnnealingFilter) deriveModuloIfPossible() {
	m.generatesEvents = false
	if m.MaximumIterations > 0 && m.percentileToReport > 0 {
		m.generatesEvents = true
		if m.percentileToReport > 1 {
			m.percentileToReport = 1
		}
		if m.percentileToReport == 1 {
			m.iterationModulo = 1
		} else {
			fmt.Printf("maximum iterations = %d\n", m.MaximumIterations)
			fmt.Printf("percentile to report = %f\n", m.percentileToReport)

			m.iterationModulo = uint64((float64)(m.MaximumIterations) * m.percentileToReport)
		}
		fmt.Printf("iteration modulo = %d\n", m.iterationModulo)

	}
}

// ShouldFilter filters only FinishedIteration Event instances, and fully filters out all StartedIteration
// events. Every FinishedIteration events received on the specified percentile boundary, one event is allowed through
// to the LogHandler.

const (
	filtered    = true
	notFiltered = false
)

func isGenerallyFilterable(eventType observer.EventType) bool {
	switch eventType {
	case observer.StartedIteration,
		observer.FinishedIteration,
		observer.Note:
		return filtered
	default:
		return notFiltered
	}
}

func isModuloFilterable(eventType observer.EventType) bool {
	switch eventType {
	case observer.StartedIteration, observer.FinishedIteration,
		observer.Note:
		return filtered
	default:
		return notFiltered
	}
}

func (m *PercentileOfIterationsPerAnnealingFilter) ShouldFilter(event observer.Event) bool {
	if !isGenerallyFilterable(event.EventType) {
		return allowThroughFilter
	}

	currentIteration := event.Attribute("CurrentIteration").(uint64)

	if m.generatesEvents && isModuloFilterable(event.EventType) &&
		currentIteration%m.iterationModulo == 0 {
		return allowThroughFilter
	}

	return blockAtFilter
}
