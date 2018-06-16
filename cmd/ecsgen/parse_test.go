package main

import (
	"bytes"
	"testing"
)

func assertCom(t *testing.T, actual, expected []Component) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("Expected to find %d component(s), but found %d", len(expected), len(actual))
		return
	}
	for i := range expected {
		if expected[i].TypeName != actual[i].TypeName {
			t.Errorf("Expected component[%d] to be named \"%s\", but found \"%s\"", i, expected[i].TypeName, actual[i].TypeName)
		}
		if expected[i].Line != actual[i].Line {
			t.Errorf("Expected %s to be on line %d, but found on line %d", expected[i].TypeName, expected[i].Line, actual[i].Line)
		}
		if expected[i].HasSerializer != actual[i].HasSerializer {
			str := "have"
			if !expected[i].HasSerializer {
				str = "not " + str
			}
			t.Errorf("Expected %s to %s a serializer", expected[i].TypeName, str)
		}
		if expected[i].HasDeserializer != actual[i].HasDeserializer {
			str := "have"
			if !expected[i].HasDeserializer {
				str = "not " + str
			}
			t.Errorf("Expected %s to %s a deserializer", expected[i].TypeName, str)
		}
	}
}

// ================
// No Serialization
// ================

const oneComNoSerialization = `package components

import (
	"github.com/danbrakeley/ecs"
)

// Fake generator line
//foo:bar ecsgen --package $GOPACKAGE --file "$GOFILE"

// PosCom is a position component
type PosCom struct {
	ecs.ComponentBase
	X      float64
	Y      float64
	Parent PosComRef
}
`

func Test_GetComponents(t *testing.T) {
	coms := getComponents(bytes.NewReader([]byte(oneComNoSerialization)))
	assertCom(t, coms, []Component{
		Component{TypeName: "PosCom", Line: 11, HasSerializer: false, HasDeserializer: false},
	})
}

// =================
//  Serializer Only
// =================

const twoComsOneSerializer = `package components

import (
	"fmt"

	"github.com/danbrakeley/ecs"
)

// Fake generator line
//foo:bar ecsgen --package $GOPACKAGE --file "$GOFILE"

// PosCom is a position component
type PosCom struct {
	ecs.ComponentBase
	X      float64
	Y      float64
	Parent PosComRef
}

func (c *PosCom) Serialize(buf *bytes.Buffer) error {
	return fmt.Errorf("not implemented")
}

// SecondCom is some other component
type SecondCom struct {
	ecs.ComponentBase
	N interface{}
}
`

func Test_GetComponents_TwoComsOneSerializer(t *testing.T) {
	coms := getComponents(bytes.NewReader([]byte(twoComsOneSerializer)))
	assertCom(t, coms, []Component{
		Component{TypeName: "PosCom", Line: 13, HasSerializer: true, HasDeserializer: false},
		Component{TypeName: "SecondCom", Line: 25, HasSerializer: false, HasDeserializer: false},
	})
}

// ===================
//  Deserializer Only
// ===================

const twoComsOneDeserializer = `package components

import (
	"fmt"

	"github.com/danbrakeley/ecs"
)

// Fake generator line
//foo:bar ecsgen --package $GOPACKAGE --file "$GOFILE"

// PosCom is a position component
type PosCom struct {
	ecs.ComponentBase
	X      float64
	Y      float64
	Parent PosComRef
}

func DeserializePosCom(mgr *ecs.Manager, buf *bytes.Buffer) error {
	return fmt.Errorf("not implemented")
}

// SecondCom is some other component
type SecondCom struct {
	ecs.ComponentBase
	N interface{}
}
`

func Test_GetComponents_TwoComsOneDeserializer(t *testing.T) {
	coms := getComponents(bytes.NewReader([]byte(twoComsOneDeserializer)))
	assertCom(t, coms, []Component{
		Component{TypeName: "PosCom", Line: 13, HasSerializer: false, HasDeserializer: true},
		Component{TypeName: "SecondCom", Line: 25, HasSerializer: false, HasDeserializer: false},
	})
}

// =============================
//  Serializer and Deserializer
// =============================

const twoComsOneBi = `package components

import (
	"fmt"

	"github.com/danbrakeley/ecs"
)

// Fake generator line
//foo:bar ecsgen --package $GOPACKAGE --file "$GOFILE"

// PosCom is a position component
type PosCom struct {
	ecs.ComponentBase
	X      float64
	Y      float64
	Parent PosComRef
}

// SecondCom is some other component
type SecondCom struct {
	ecs.ComponentBase
	N interface{}
}

func (c *SecondCom) Serialize(buf *bytes.Buffer) error {
	return fmt.Errorf("not implemented")
}

func DeserializeSecondCom(mgr *ecs.Manager, buf *bytes.Buffer) error {
	return fmt.Errorf("not implemented")
}
`

func Test_GetComponents_TwoComsOneBiDirectional(t *testing.T) {
	coms := getComponents(bytes.NewReader([]byte(twoComsOneBi)))
	assertCom(t, coms, []Component{
		Component{TypeName: "PosCom", Line: 13, HasSerializer: false, HasDeserializer: false},
		Component{TypeName: "SecondCom", Line: 21, HasSerializer: true, HasDeserializer: true},
	})
}
