package ecs

import "bytes"

// ComponentName provides a unique type for names of components
type ComponentName string

// Component interface
type Component interface {
	GetName() ComponentName
	GetEntity() *Entity
	SetEntity(e *Entity)
	Serialize(buf *bytes.Buffer) error
}

// ComponentBase can be embedded in a Component to provide owning Entity tracking
type ComponentBase struct {
	entity *Entity
}

// GetEntity returns the owning Entity of this Component
func (c *ComponentBase) GetEntity() *Entity {
	return c.entity
}

// SetEntity sets the Entity that owns this Component
func (c *ComponentBase) SetEntity(e *Entity) {
	c.entity = e
}

// GetEntityID returns the owning Entity's ID
func (c *ComponentBase) GetEntityID() string {
	return c.entity.GetID()
}

// GetManager returns the owning Entity's Manager
func (c *ComponentBase) GetManager() *Manager {
	return c.entity.GetManager()
}
