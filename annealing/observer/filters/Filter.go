// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// modulators package supplies a number of observer modulators for managing the chattiness of loggers.
package filters

import "github.com/LindsayBradford/crem/annealing"

// Filter describes an interface to object that decides on how observer should be filtered
// (reduced in volume of entries logged).
type Filter interface {

	// ShouldFilter accepts an Event and decides whether it should be filtered (not logged).
	// This method returns true iff the logger is to ignore observer the supplied event.
	ShouldFilter(event annealing.Event) bool
}

// NullFilter is a default Filter that doesn't actually filter observer, allowing all events through.
type NullFilter struct{}

// ShouldFilter always returns false (do not filter the log of events)
func (nm *NullFilter) ShouldFilter(event annealing.Event) bool {
	return false
}