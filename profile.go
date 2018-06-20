package ecs

import (
	"fmt"
	"time"
)

// Stats is a collection of metrics
type Stats struct {
	ManagerCreate   time.Time
	NumEntityCreate int64
	NumEntityDelete int64
	Sys             map[string]UpdateStats
}

// NewStats is the Stats constructor
func NewStats(t time.Time) Stats {
	return Stats{
		ManagerCreate: t,
		Sys:           make(map[string]UpdateStats),
	}
}

// Print dumps the current stats to stdout
func (s Stats) Print() {
	now := time.Now()
	fmt.Printf("--- STATS at %v\n", now)
	fmt.Printf("Running Time: %v\n", now.Sub(s.ManagerCreate))
	fmt.Printf("Entity Create/Delete: %d/%d\n", s.NumEntityCreate, s.NumEntityDelete)
	for k, v := range s.Sys {
		avg := v.TotalTime / time.Duration(v.TotalFrames)
		fmt.Printf("%s: Avg: %v, Best: %v, Worst: %v\n", k, avg, v.BestFrame, v.WorstFrame)
	}
}

// AddSystemFrame updates or creates stats for a given system
func (s *Stats) AddSystemFrame(Name string, dt time.Duration) {
	u, ok := s.Sys[Name]
	if !ok {
		u = UpdateStats{BestFrame: time.Hour}
	}
	u.AddFrame(dt)
	s.Sys[Name] = u
}

// UpdateStats is a collection of metrics about an individual System's FixedUpdate
type UpdateStats struct {
	TotalTime   time.Duration
	TotalFrames int64
	WorstFrame  time.Duration
	BestFrame   time.Duration
}

// AddFrame updates stats based on the given dt
func (s *UpdateStats) AddFrame(dt time.Duration) {
	s.TotalFrames++
	s.TotalTime += dt
	if dt < s.BestFrame {
		s.BestFrame = dt
	}
	if dt > s.WorstFrame {
		s.WorstFrame = dt
	}
}
