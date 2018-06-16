package ecs

import (
	"bytes"
	"fmt"
	"time"

	"github.com/teris-io/shortid"
)

// Manager is the glue between a set of systems and their components.
// Manager also has a World, which groups entites and allows serializable
// references between components of different entities.
type Manager struct {
	world     *World
	systems   []System
	updaters  []FixedUpdater
	tickers   []Ticker
	receivers []Receiver
}

// NewManager constructs a new Manager instance
func NewManager() *Manager {
	m := Manager{}
	m.world = NewWorld(&m)
	return &m
}

// MustRegisterSystem calls RegisterSystem and panics if there's an error
func (m *Manager) MustRegisterSystem(s System) {
	err := m.RegisterSystem(s)
	if err != nil {
		panic(err)
	}
}

// RegisterSystem adds a system to the Manager.
// It should be called before any components are added to entities.
func (m *Manager) RegisterSystem(s System) error {
	name := s.GetName()
	for _, v := range m.systems {
		if v.GetName() == name {
			return fmt.Errorf("cannot register system %s because it is already registered", name)
		}
	}
	m.systems = append(m.systems, s)

	if updater, ok := s.(FixedUpdater); ok {
		m.updaters = append(m.updaters, updater)
	}
	if ticker, ok := s.(Ticker); ok {
		m.tickers = append(m.tickers, ticker)
	}
	if receiver, ok := s.(Receiver); ok {
		m.receivers = append(m.receivers, receiver)
	}

	return nil
}

// FixedUpdate calls FixedUpdate on each registered system
func (m *Manager) FixedUpdate(dt float64) {
	for _, s := range m.updaters {
		s.FixedUpdate(dt)
	}
}

// Tick calls Tick on each registered system
func (m *Manager) Tick(dt time.Duration) {
	for _, t := range m.tickers {
		t.Tick(dt)
	}
}

// NewEntity constructs a new Entity
func (m *Manager) NewEntity() *Entity {
	return m.newEntityWithID(shortid.MustGenerate())
}

// RecreateEntity constructs an Entity with a known ID
func (m *Manager) RecreateEntity(id string) *Entity {
	return m.newEntityWithID(id)
}

// newEntityWithID constructs an entity with the given ID
func (m *Manager) newEntityWithID(id string) *Entity {
	e := &Entity{
		manager:    m,
		components: make(map[ComponentName]Component),
		id:         id,
		isDead:     false,
	}
	m.world.Add(e)
	return e
}

// GetEntity retreives an entity from its id
func (m *Manager) GetEntity(id string) *Entity {
	return m.world.Find(id)
}

// UpdateEntity is an internal func, meant to be called by an Entity when its components have changed.
// This gives all systems a chance to decide if they care about changes to the composition of components.
func (m *Manager) UpdateEntity(e *Entity) {
	for _, s := range m.systems {
		s.UpdateEntity(e)
	}
}

// DropEntity is an internal func, meant to be called by an Entity when it wants to be removed.
// All systems should remove this entity and all components from their internal lists.
func (m *Manager) DropEntity(e *Entity) {
	for _, s := range m.systems {
		s.DropEntity(e)
	}
	m.world.Remove(e)
}

// Broadcast sends the given entity to all systems that implement the Receiver interface
func (m *Manager) Broadcast(e *Entity) {
	for _, r := range m.receivers {
		r.HandleBroadcast(e)
	}
}

// Serialize saves all entities (and their components)
func (m *Manager) Serialize(buf *bytes.Buffer) error {
	return m.world.Serialize(buf)
}

// Deserialize loads entities (and their components)
func (m *Manager) Deserialize(buf *bytes.Buffer) error {
	return m.world.Deserialize(buf)
}
