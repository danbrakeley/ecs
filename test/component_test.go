package test

import (
	"bytes"
	"testing"

	"github.com/danbrakeley/ecs"
)

//go:generate ecsgen --package "$GOPACKAGE" --file "$GOFILE"

// TestEmptyCom is a component with no fields
type TestEmptyCom struct {
	ecs.ComponentBase
}

// DeserializeTestEmptyCom is this component's deserialization function
func DeserializeTestEmptyCom(mgr *ecs.Manager, e *ecs.Entity, buf *bytes.Buffer) error {
	return nil
}

//
// Tests

func Test_Component_ComponentSetsOwningEntity(t *testing.T) {
	mgr, _ := createMgrAndBaseSys()
	c := &TestEmptyCom{}
	e := mgr.NewEntity()
	e.AddComponent(c)

	if c.GetEntity() == nil || e.GetID() != c.GetEntity().GetID() {
		t.Errorf("Added component did not set owning entity")
	}
}

func Test_Component_CannotAddComponentToTwoEntities(t *testing.T) {
	mgr, _ := createMgrAndBaseSys()
	c := &TestEmptyCom{}
	e1 := mgr.NewEntity()
	e1.AddComponent(c)
	e2 := mgr.NewEntity()

	if !didPanic(func() { e2.AddComponent(c) }) {
		t.Errorf("Added component to two different entities without a panic")
	}
}

func Test_Component_ReferenceIsNil(t *testing.T) {
	mgr := ecs.NewManager()

	e1 := mgr.NewEntity()
	c1 := TestEmptyCom{}
	e1.AddComponent(&c1)

	e2 := mgr.NewEntity()
	c2 := TestSerialRefCom{}
	e2.AddComponent(&c2)

	if !c2.Ref.IsNil() {
		t.Errorf("Component reference is not nil (but should be)")
	}

	c2.Ref = NewTestEmptyComRef(&c1)

	if c2.Ref.IsNil() {
		t.Errorf("Component reference is nil (but shouldn't be)")
	}
}

func Test_Component_ReferenceGet(t *testing.T) {
	mgr := ecs.NewManager()

	e1 := mgr.NewEntity()
	c1 := TestEmptyCom{}
	e1.AddComponent(&c1)

	e2 := mgr.NewEntity()
	c2 := TestSerialRefCom{Ref: NewTestEmptyComRef(&c1)}
	e2.AddComponent(&c2)

	if c2.Ref.Get() != &c1 {
		t.Errorf("Component reference did not resolve")
	}
}
