package server

import (
	"sync"
)

// Scheduler struct to manage the scheduling system
type Scheduler struct {
	loaded   map[string]*runnerRef
	loadedMu sync.Mutex
}

// InitScheduler initializes the scheduling system for the application.
func InitScheduler() *Scheduler {
	return &Scheduler{
		loaded: make(map[string]*runnerRef),
	}
}
