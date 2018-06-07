package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/danbrakeley/ecs"
)

//go:generate ecsgen --package "$GOPACKAGE" --file "$GOFILE"

// TestSerialCom is a serializable component
type TestSerialCom struct {
	ecs.ComponentBase
	N int `json:"n"`
}

func (c *TestSerialCom) Serialize(buf *bytes.Buffer) error {
	if err := ecs.SerializeToJSON(buf, c); err != nil {
		return err
	}
	return nil
}

// DeserializeTestSerialCom creates a TestSerialCom from the given buffer
func DeserializeTestSerialCom(mgr *ecs.Manager, e *ecs.Entity, buf *bytes.Buffer) error {
	var com TestSerialCom
	if err := json.Unmarshal(buf.Bytes(), &com); err != nil {
		return err
	}
	e.AddComponent(&com)
	return nil
}

func Test_Serialize_SingleField(t *testing.T) {
	mgr := ecs.NewManager()
	e := mgr.NewEntity()
	e.AddComponent(&TestSerialCom{N: 1})

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	err := e.Serialize(buf)
	if err != nil {
		t.Fatalf("Error in Serialize: %v", err)
	}

	actual := string(buf.Bytes())
	expected := fmt.Sprintf(`{"id":"%s","components":{"TestSerialCom":{"n":1}}}`, e.GetID())
	if diff := formatStringDiff(expected, actual); diff != nil {
		t.Fatalf("Serialized bytes mismatch:\n%s", strings.Join(diff, "\n"))
	}
}

func Test_Deserialize_SingleEntity(t *testing.T) {
	mgr := ecs.NewManager()
	e := mgr.NewEntity()
	e.AddComponent(&TestSerialCom{N: 1})

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	err := e.Serialize(buf)
	if err != nil {
		t.Fatalf("Error in Serialize: %v", err)
	}

	mgr = ecs.NewManager()
	eNew, err := ecs.DeserializeEntity(mgr, buf)
	if err != nil {
		t.Fatalf("Error in DeserializeEntity: %v", err)
	}
	cNew := GetTestSerialCom(eNew)
	if cNew == nil {
		t.Fatalf("Deserialized entity missing TestSerialCom")
	}

	if cNew.N != 1 {
		t.Errorf("Deserialized TestSerialCom has %d for N, expected 1", cNew.N)
	}
}

// TestSerialMultiCom is a serializable component with multiple fields of different types
type TestSerialMultiCom struct {
	ecs.ComponentBase
	Foo    string `json:"foo"`
	BarBaz string `json:"bar_baz"`
}

func (c *TestSerialMultiCom) Serialize(buf *bytes.Buffer) error {
	if err := ecs.SerializeToJSON(buf, c); err != nil {
		return err
	}
	return nil
}

func (c *TestSerialMultiCom) Deserialize(buf *bytes.Buffer) error {
	return nil
}

// TestSerialRefCom is a serializable component with a reference to another entity's component
type TestSerialRefCom struct {
	ecs.ComponentBase
	Ref TestEmptyComRef `json:"ref"`
}

func (c *TestSerialRefCom) Serialize(buf *bytes.Buffer) error {
	buf.WriteString(`{"ref":`)
	c.Ref.Serialize(buf)
	buf.WriteString(`}`)
	return nil
}

func (c *TestSerialRefCom) Deserialize(buf *bytes.Buffer) error {
	return nil
}

//
// Serialize Tests

func Test_Serialize_MultipleFields(t *testing.T) {
	mgr := ecs.NewManager()
	e := mgr.NewEntity()
	e.AddComponent(&TestSerialMultiCom{Foo: "a", BarBaz: "b"})

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	err := e.Serialize(buf)
	if err != nil {
		t.Fatalf("Error in Serialize: %v", err)
	}

	actual := string(buf.Bytes())
	expected := fmt.Sprintf(`{"id":"%s","components":{"TestSerialMultiCom":{"foo":"a","bar_baz":"b"}}}`, e.GetID())
	if diff := formatStringDiff(expected, actual); diff != nil {
		t.Fatalf("Serialized bytes:\n%s", strings.Join(diff, "\n"))
	}
}

func Test_Serialize_Reference(t *testing.T) {
	mgr := ecs.NewManager()
	e1 := mgr.NewEntity()
	sc := TestEmptyCom{}
	e1.AddComponent(&sc)

	e2 := mgr.NewEntity()
	rc := TestSerialRefCom{Ref: NewTestEmptyComRef(&sc)}
	e2.AddComponent(&rc)

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	err := e2.Serialize(buf)
	if err != nil {
		t.Fatalf("Error in Serialize: %v", err)
	}

	actual := string(buf.Bytes())
	expected := fmt.Sprintf(`{"id":"%s","components":{"TestSerialRefCom":{"ref":"%s"}}}`, e2.GetID(), e1.GetID())
	if diff := formatStringDiff(expected, actual); diff != nil {
		t.Fatalf("Serialized bytes:\n%s", strings.Join(diff, "\n"))
	}
}
