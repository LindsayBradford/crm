// Copyright (c) 2018 Australian Rivers Institute.

// Package rand implements an concurrency-safe wrapper around the language supplied rand package for
// pseudo-random number generators.  It does so by wrapping calls to an underlying math/rand Rand struct
// with a mutex lock.
package rand

import (
	"math/rand"
	"sync"
	"time"
)

// ConcurrencySafeRand is a concurrency-safe source of random numbers
type ConcurrencySafeRand struct {
	unsafeRand *rand.Rand
	mutex      sync.Mutex
}

// New returns a new ConcurrencySafeRand that uses random values from src to generate other random values.
func New(src rand.Source) *ConcurrencySafeRand {
	unsafeRand := rand.New(src)
	mutex := sync.Mutex{}
	return &ConcurrencySafeRand{unsafeRand: unsafeRand, mutex: mutex}
}

// New returns a new ConcurrencySafeRand that uses random values seeded from a source of the system-time to generate
// other random values.
func NewTimeSeeded() *ConcurrencySafeRand {
	return New(rand.NewSource(time.Now().UnixNano()))
}

// Uint64 returns a pseudo-random 64-bit value as a uint64 from the default Source.
func (csr *ConcurrencySafeRand) Uint64() uint64 {
	csr.mutex.Lock()
	defer csr.mutex.Unlock()
	return csr.unsafeRand.Uint64()
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (csr *ConcurrencySafeRand) Intn(n int) int {
	csr.mutex.Lock()
	defer csr.mutex.Unlock()
	return csr.unsafeRand.Intn(n)
}
