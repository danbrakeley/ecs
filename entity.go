package ecs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Entity is a container for the components that make up a single object
type Entity struct {
	manager    *Manager
	components map[ComponentName]Component
	id         string
	isDead     bool // only true if this entity has been Dropped
}

// GetManager returns this entity's owning manager
func (e *Entity) GetManager() *Manager {
	return e.manager
}

// AddComponent adds the component instance to this Entity.
// This call also updates all Systems that might care.
// Returns *Entity to allow for chaining of other calls.
func (e *Entity) AddComponent(c Component) *Entity {
	if e.isDead {
		panic(fmt.Errorf("Cannot add %s to entity %s because the entity is dead", c.GetName(), e.GetID()))
	}
	oldEntity := c.GetEntity()
	if oldEntity != nil {
		panic(fmt.Errorf("Cannot add %s to entity %s because it already belongs to entity %s",
			c.GetName(), e.GetID(), oldEntity.GetID()),
		)
	}
	c.SetEntity(e)
	e.components[c.GetName()] = c
	e.manager.UpdateEntity(e)
	return e
}

// RemoveComponent removes the component instance from this Entity.
// This call also updates all Systems that might care.
// Returns *Entity to allow for chaining of other calls.
func (e *Entity) RemoveComponent(c Component) *Entity {
	if e.isDead {
		panic(fmt.Errorf("Cannot remove %s from entity %s because the entity is dead", c.GetName(), e.GetID()))
	}
	oldEntity := c.GetEntity()
	if oldEntity == nil || oldEntity.GetID() != e.GetID() {
		panic(fmt.Errorf("Component %s doesn't belong to entity %s, and thus cannot be removed", c.GetName(), e.GetID()))
	}
	c.SetEntity(nil)
	delete(e.components, c.GetName())
	e.manager.UpdateEntity(e)
	return e
}

// GetComponent returns the named Component, or nil if not found
func (e *Entity) GetComponent(name ComponentName) Component {
	return e.components[name]
}

// GetID returns this Entity's unique ID
func (e *Entity) GetID() string {
	return e.id
}

// Drop removes this entity's components from all systems.
// This entity and its components will be GC'd once all references are gone.
// The dropped components will report nil from GetEntity().
func (e *Entity) Drop() {
	if e.isDead {
		return
	}
	// tell systems to drop the entity
	e.manager.DropEntity(e)
	// tell all components to drop their owner
	for _, c := range e.components {
		myID := e.GetID()
		comOwnerID := c.GetEntity().GetID()
		if myID != comOwnerID {
			panic(fmt.Errorf("Entity %s has %s which thinks it is owned by entity %s", myID, c.GetName(), comOwnerID))
		}
		c.SetEntity(nil)
	}
	e.isDead = true
}

// Serialize saves this entity and all its serializable components
func (e *Entity) Serialize(buf *bytes.Buffer) error {
	if e.isDead {
		return fmt.Errorf("tried to serialize dead entity %s", e.GetID())
	}
	buf.WriteString("{\"id\":\"")
	buf.WriteString(e.GetID())
	buf.WriteString("\",\"components\":{")

	comOrder := make([]string, len(e.components))
	i := 0
	for comName := range e.components {
		comOrder[i] = string(comName)
		i++
	}
	sort.Strings(comOrder)

	for i, comName := range comOrder {
		if i != 0 {
			buf.WriteRune(',')
		}
		com := e.components[ComponentName(comName)]
		buf.WriteString(fmt.Sprintf(`"%s":`, comName))
		if err := com.Serialize(buf); err != nil {
			return err
		}
	}

	buf.WriteString("}}")
	return nil
}

type entityRaw struct {
	ID         string                            `json:"id"`
	Components map[ComponentName]json.RawMessage `json:"components"`
}

// DeserializeEntity restores an entity (and its components)
func DeserializeEntity(mgr *Manager, buf *bytes.Buffer) (*Entity, error) {
	var entity entityRaw
	if err := json.Unmarshal(buf.Bytes(), &entity); err != nil {
		return nil, err
	}

	e := mgr.RecreateEntity(entity.ID)
	for comName, rawBytes := range entity.Components {
		fnDeserialize := GetFnDeserializer(comName)
		if fnDeserialize == nil {
			return nil, fmt.Errorf("No Deserializer registered for %s", comName)
		}
		if err := fnDeserialize(mgr, e, bytes.NewBuffer(rawBytes)); err != nil {
			return nil, err
		}
	}

	return e, nil
}
