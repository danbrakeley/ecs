package test

import (
	"bytes"
	"fmt"

	"github.com/danbrakeley/ecs"
)

//   ___  ___ ___  __ _  ___ _ __
//  / _ \/ __/ __|/ _` |/ _ \ '_ \
// |  __/ (__\__ \ (_| |  __/ | | |  v0.3.0
//  \___|\___|___/\__, |\___|_| |_|  2018-06-06T22:47:02-04:00
//                |___/
//
// WARNING: This file was generated by ecsgen.
// Any changes made to this file by hand may be lost.
//

func init() {
	ecs.RegisterComponent(TestEmptyComName, DeserializeTestEmptyCom)
}

//
// *** TestEmptyCom ***

// TestEmptyComName is TestEmptyCom's ComponentName
const TestEmptyComName ecs.ComponentName = "TestEmptyCom"

// GetTestEmptyCom returns any TestEmptyCom on the given Entity
func GetTestEmptyCom(e *ecs.Entity) *TestEmptyCom {
	if c := e.GetComponent(TestEmptyComName); c != nil {
		return c.(*TestEmptyCom)
	}
	return nil
}

// GetName is from the Component interface
func (c *TestEmptyCom) GetName() ecs.ComponentName {
	return TestEmptyComName
}

// Serialize writes this TestEmptyCom to the given buffer
func (c *TestEmptyCom) Serialize(buf *bytes.Buffer) error {
	return ecs.SerializeToJSON(buf, c)
}

//
// TestEmptyComRef is a TestEmptyCom reference that can be serialized
type TestEmptyComRef struct {
	mgr      *ecs.Manager
	parentID string
}

// NewTestEmptyComRef constructs a TestEmptyComRef
func NewTestEmptyComRef(c *TestEmptyCom) TestEmptyComRef {
	var r TestEmptyComRef
	r.Set(c)
	return r
}

// Set updates this reference to point to the given component (or nil)
func (r *TestEmptyComRef) Set(c *TestEmptyCom) {
	if c == nil {
		r.mgr = nil
		r.parentID = ""
		return
	}
	r.mgr = c.GetManager()
	r.parentID = c.GetEntityID()
}

// IsNil checks if the component pointer == nil
func (r TestEmptyComRef) IsNil() bool {
	return len(r.parentID) == 0
}

// Get resolves this reference to a TestEmptyCom pointer
func (r TestEmptyComRef) Get() *TestEmptyCom {
	if r.IsNil() {
		return nil
	}
	return GetTestEmptyCom(r.mgr.GetEntity(r.parentID))
}

// GetEntity resolves the Entity owning the referenced component
func (r TestEmptyComRef) GetEntity() *ecs.Entity {
	if r.IsNil() {
		return nil
	}
	return r.mgr.GetEntity(r.parentID)
}

// Serialize just writes out the parent entity id
func (r TestEmptyComRef) Serialize(buf *bytes.Buffer) error {
	if r.IsNil() {
		buf.WriteString("null")
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\"", r.parentID))
	return nil
}