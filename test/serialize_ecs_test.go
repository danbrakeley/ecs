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
	ecs.RegisterComponent(TestSerialComName, DeserializeTestSerialCom)
	ecs.RegisterComponent(TestSerialMultiComName, DeserializeTestSerialMultiCom)
	ecs.RegisterComponent(TestSerialRefComName, DeserializeTestSerialRefCom)
}

//
// *** TestSerialCom ***

// TestSerialComName is TestSerialCom's ComponentName
const TestSerialComName ecs.ComponentName = "TestSerialCom"

// GetTestSerialCom returns any TestSerialCom on the given Entity
func GetTestSerialCom(e *ecs.Entity) *TestSerialCom {
	if c := e.GetComponent(TestSerialComName); c != nil {
		return c.(*TestSerialCom)
	}
	return nil
}

// GetName is from the Component interface
func (c *TestSerialCom) GetName() ecs.ComponentName {
	return TestSerialComName
}

//
// TestSerialComRef is a TestSerialCom reference that can be serialized
type TestSerialComRef struct {
	mgr      *ecs.Manager
	parentID string
}

// NewTestSerialComRef constructs a TestSerialComRef
func NewTestSerialComRef(c *TestSerialCom) TestSerialComRef {
	var r TestSerialComRef
	r.Set(c)
	return r
}

// Set updates this reference to point to the given component (or nil)
func (r *TestSerialComRef) Set(c *TestSerialCom) {
	if c == nil {
		r.mgr = nil
		r.parentID = ""
		return
	}
	r.mgr = c.GetManager()
	r.parentID = c.GetEntityID()
}

// IsNil checks if the component pointer == nil
func (r TestSerialComRef) IsNil() bool {
	return len(r.parentID) == 0
}

// Get resolves this reference to a TestSerialCom pointer
func (r TestSerialComRef) Get() *TestSerialCom {
	if r.IsNil() {
		return nil
	}
	return GetTestSerialCom(r.mgr.GetEntity(r.parentID))
}

// GetEntity resolves the Entity owning the referenced component
func (r TestSerialComRef) GetEntity() *ecs.Entity {
	if r.IsNil() {
		return nil
	}
	return r.mgr.GetEntity(r.parentID)
}

// Serialize just writes out the parent entity id
func (r TestSerialComRef) Serialize(buf *bytes.Buffer) error {
	if r.IsNil() {
		buf.WriteString("null")
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\"", r.parentID))
	return nil
}

//
// *** TestSerialMultiCom ***

// TestSerialMultiComName is TestSerialMultiCom's ComponentName
const TestSerialMultiComName ecs.ComponentName = "TestSerialMultiCom"

// GetTestSerialMultiCom returns any TestSerialMultiCom on the given Entity
func GetTestSerialMultiCom(e *ecs.Entity) *TestSerialMultiCom {
	if c := e.GetComponent(TestSerialMultiComName); c != nil {
		return c.(*TestSerialMultiCom)
	}
	return nil
}

// GetName is from the Component interface
func (c *TestSerialMultiCom) GetName() ecs.ComponentName {
	return TestSerialMultiComName
}

// DeserializeTestSerialMultiCom creates a TestSerialMultiCom from the given buffer
func DeserializeTestSerialMultiCom(mgr *ecs.Manager, e *ecs.Entity, buf *bytes.Buffer) error {
	return fmt.Errorf("DeserializeTestSerialMultiCom not implemented")
}

//
// TestSerialMultiComRef is a TestSerialMultiCom reference that can be serialized
type TestSerialMultiComRef struct {
	mgr      *ecs.Manager
	parentID string
}

// NewTestSerialMultiComRef constructs a TestSerialMultiComRef
func NewTestSerialMultiComRef(c *TestSerialMultiCom) TestSerialMultiComRef {
	var r TestSerialMultiComRef
	r.Set(c)
	return r
}

// Set updates this reference to point to the given component (or nil)
func (r *TestSerialMultiComRef) Set(c *TestSerialMultiCom) {
	if c == nil {
		r.mgr = nil
		r.parentID = ""
		return
	}
	r.mgr = c.GetManager()
	r.parentID = c.GetEntityID()
}

// IsNil checks if the component pointer == nil
func (r TestSerialMultiComRef) IsNil() bool {
	return len(r.parentID) == 0
}

// Get resolves this reference to a TestSerialMultiCom pointer
func (r TestSerialMultiComRef) Get() *TestSerialMultiCom {
	if r.IsNil() {
		return nil
	}
	return GetTestSerialMultiCom(r.mgr.GetEntity(r.parentID))
}

// GetEntity resolves the Entity owning the referenced component
func (r TestSerialMultiComRef) GetEntity() *ecs.Entity {
	if r.IsNil() {
		return nil
	}
	return r.mgr.GetEntity(r.parentID)
}

// Serialize just writes out the parent entity id
func (r TestSerialMultiComRef) Serialize(buf *bytes.Buffer) error {
	if r.IsNil() {
		buf.WriteString("null")
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\"", r.parentID))
	return nil
}

//
// *** TestSerialRefCom ***

// TestSerialRefComName is TestSerialRefCom's ComponentName
const TestSerialRefComName ecs.ComponentName = "TestSerialRefCom"

// GetTestSerialRefCom returns any TestSerialRefCom on the given Entity
func GetTestSerialRefCom(e *ecs.Entity) *TestSerialRefCom {
	if c := e.GetComponent(TestSerialRefComName); c != nil {
		return c.(*TestSerialRefCom)
	}
	return nil
}

// GetName is from the Component interface
func (c *TestSerialRefCom) GetName() ecs.ComponentName {
	return TestSerialRefComName
}

// DeserializeTestSerialRefCom creates a TestSerialRefCom from the given buffer
func DeserializeTestSerialRefCom(mgr *ecs.Manager, e *ecs.Entity, buf *bytes.Buffer) error {
	return fmt.Errorf("DeserializeTestSerialRefCom not implemented")
}

//
// TestSerialRefComRef is a TestSerialRefCom reference that can be serialized
type TestSerialRefComRef struct {
	mgr      *ecs.Manager
	parentID string
}

// NewTestSerialRefComRef constructs a TestSerialRefComRef
func NewTestSerialRefComRef(c *TestSerialRefCom) TestSerialRefComRef {
	var r TestSerialRefComRef
	r.Set(c)
	return r
}

// Set updates this reference to point to the given component (or nil)
func (r *TestSerialRefComRef) Set(c *TestSerialRefCom) {
	if c == nil {
		r.mgr = nil
		r.parentID = ""
		return
	}
	r.mgr = c.GetManager()
	r.parentID = c.GetEntityID()
}

// IsNil checks if the component pointer == nil
func (r TestSerialRefComRef) IsNil() bool {
	return len(r.parentID) == 0
}

// Get resolves this reference to a TestSerialRefCom pointer
func (r TestSerialRefComRef) Get() *TestSerialRefCom {
	if r.IsNil() {
		return nil
	}
	return GetTestSerialRefCom(r.mgr.GetEntity(r.parentID))
}

// GetEntity resolves the Entity owning the referenced component
func (r TestSerialRefComRef) GetEntity() *ecs.Entity {
	if r.IsNil() {
		return nil
	}
	return r.mgr.GetEntity(r.parentID)
}

// Serialize just writes out the parent entity id
func (r TestSerialRefComRef) Serialize(buf *bytes.Buffer) error {
	if r.IsNil() {
		buf.WriteString("null")
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\"", r.parentID))
	return nil
}