package ecs

import (
	"bytes"
)

// FnDeserialize is the type of a Component's Deserialization routine
// It recreates its Component in the given Entity, reading from buf, and using the given Manager
type FnDeserialize func(mgr *Manager, e *Entity, buf *bytes.Buffer) error

var (
	comRegsitry map[ComponentName]FnDeserialize
)

func init() {
	comRegsitry = make(map[ComponentName]FnDeserialize)
}

// RegisterComponent must be called for a component to be able to be deserialized
func RegisterComponent(comName ComponentName, fn FnDeserialize) {
	comRegsitry[comName] = fn
}

// GetFnDeserializer returns the FnDeserializer for the given component, or nil
func GetFnDeserializer(comName ComponentName) FnDeserialize {
	return comRegsitry[comName]
}
