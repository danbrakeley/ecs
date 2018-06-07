package ecs

import (
	"bytes"
	"encoding/json"
	"sort"
)

// World is a container for entities.
// Currently, the main use is to preserve inter-entity references during Deserialization.
type World struct {
	mgr      *Manager
	entities map[string]*Entity
}

// NewWorld constructs an World
func NewWorld(mgr *Manager) *World {
	return &World{
		mgr:      mgr,
		entities: make(map[string]*Entity),
	}
}

// Add adds a unique entity to this container.
// If there is already an entry for this entity's ID, this new entity replaces the existing one.
func (w *World) Add(e *Entity) {
	w.entities[e.GetID()] = e
}

// Remove an entity from this container.
func (w *World) Remove(e *Entity) {
	delete(w.entities, e.GetID())
}

// Find looks up the entity pointer from its ID
// TODO: decide if I like Get or Find better
func (w *World) Find(id string) *Entity {
	return w.entities[id]
}

// Get looks up the entity pointer from its ID
// TODO: decide if I like Get or Find better
func (w *World) Get(id string) *Entity {
	return w.entities[id]
}

// Serialize saves all entities (and their components)
func (w *World) Serialize(buf *bytes.Buffer) error {
	buf.WriteString("{\"entities\":[")

	// ensure consistent ordering in serialized data
	keys := make([]string, 0, len(w.entities))
	for k := range w.entities {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// serialize in order
	for i, k := range keys {
		if i != 0 {
			buf.WriteRune(',')
		}
		if err := w.entities[k].Serialize(buf); err != nil {
			return err
		}
	}

	buf.WriteString("]}")
	return nil
}

type worldRaw struct {
	Entities []json.RawMessage `json:"entities"`
}

// Deserialize restores all entities (and their components)
func (w *World) Deserialize(buf *bytes.Buffer) error {
	// TODO: probably more efficient to deserialize buffer in one pass, rather than
	// deserializing one buffer into multiple smaller buffers like I do here.
	var world worldRaw
	err := json.Unmarshal(buf.Bytes(), &world)
	if err != nil {
		return err
	}

	for _, raw := range world.Entities {
		e, err := DeserializeEntity(w.mgr, bytes.NewBuffer([]byte(raw)))
		if err != nil {
			return err
		}
		id := e.GetID()

		// TODO: probably more efficient way to swap in a new entity than just deleting the old one
		if ePrev, ok := w.entities[id]; ok {
			// delete any previous entity
			w.mgr.DropEntity(ePrev)
		}
		w.entities[e.GetID()] = e
		// add new entity
		w.mgr.UpdateEntity(e)
	}

	return nil
}
