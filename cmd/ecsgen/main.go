package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
)

const (
	version         = "0.3.1"
	templatedOutput = `package {{.Package}}

import (
	"bytes"
	"fmt"

	"github.com/danbrakeley/ecs"
)

//   ___  ___ ___  __ _  ___ _ __
//  / _ \/ __/ __|/ _{{.Backtick}} |/ _ \ '_ \
// |  __/ (__\__ \ (_| |  __/ | | |
//  \___|\___|___/\__, |\___|_| |_|  v{{.Version}}
//                |___/
//
// WARNING: This file was generated by ecsgen.
// Any changes made to this file by hand may be lost.
//

func init() {
{{range .Components}}	ecs.RegisterComponent({{.Name}}Name, Deserialize{{.Name}})
{{end}}}
{{range .Components}}
//
// *** {{.Name}} ***

// {{.Name}}Name is {{.Name}}'s ComponentName
const {{.Name}}Name ecs.ComponentName = "{{.Name}}"

// Get{{.Name}} returns any {{.Name}} on the given Entity
func Get{{.Name}}(e *ecs.Entity) *{{.Name}} {
	if c := e.GetComponent({{.Name}}Name); c != nil {
		return c.(*{{.Name}})
	}
	return nil
}

// GetName is from the Component interface
func (c *{{.Name}}) GetName() ecs.ComponentName {
	return {{.Name}}Name
}
{{if (eq .HasSerializer false)}}
// Serialize writes this {{.Name}} to the given buffer
func (c *{{.Name}}) Serialize(buf *bytes.Buffer) error {
	return ecs.SerializeToJSON(buf, c)
}
{{end}}{{if (eq .HasDeserializer false)}}
// Deserialize{{.Name}} creates a {{.Name}} from the given buffer
func Deserialize{{.Name}}(mgr *ecs.Manager, e *ecs.Entity, buf *bytes.Buffer) error {
	return fmt.Errorf("Deserialize{{.Name}} not implemented")
}
{{end}}
//
// {{.Name}}Ref is a {{.Name}} reference that can be serialized
type {{.Name}}Ref struct {
	mgr      *ecs.Manager
	parentID string
}

// New{{.Name}}Ref constructs a {{.Name}}Ref
func New{{.Name}}Ref(c *{{.Name}}) {{.Name}}Ref {
	var r {{.Name}}Ref
	r.Set(c)
	return r
}

// Set updates this reference to point to the given component (or nil)
func (r *{{.Name}}Ref) Set(c *{{.Name}}) {
	if c == nil {
		r.mgr = nil
		r.parentID = ""
		return
	}
	r.mgr = c.GetManager()
	r.parentID = c.GetEntityID()
}

// IsNil checks if the component pointer == nil
func (r {{.Name}}Ref) IsNil() bool {
	return len(r.parentID) == 0
}

// Get resolves this reference to a {{.Name}} pointer
func (r {{.Name}}Ref) Get() *{{.Name}} {
	if r.IsNil() {
		return nil
	}
	return Get{{.Name}}(r.mgr.GetEntity(r.parentID))
}

// GetEntity resolves the Entity owning the referenced component
func (r {{.Name}}Ref) GetEntity() *ecs.Entity {
	if r.IsNil() {
		return nil
	}
	return r.mgr.GetEntity(r.parentID)
}

// Serialize just writes out the parent entity id
func (r {{.Name}}Ref) Serialize(buf *bytes.Buffer) error {
	if r.IsNil() {
		buf.WriteString("null")
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\"", r.parentID))
	return nil
}
{{end}}`
)

// getOutputFileName returns a string in the form: <base>_ecs[_test].ext
func getOutputFileName(inFile string) string {
	ext := path.Ext(inFile)
	base := strings.TrimSuffix(inFile, ext)
	var genFmt string
	if strings.HasSuffix(base, "_test") {
		base = strings.TrimSuffix(base, "_test")
		genFmt = "%s_ecs_test%s"
	} else {
		genFmt = "%s_ecs%s"
	}
	return fmt.Sprintf(genFmt, base, ext)
}

// Vars are the dynamic fields that are needed to generate a component _gen.go file
type Vars struct {
	Backtick   string // a string containing a single backtick (`)
	Version    string
	Package    string
	File       string
	GenFile    string
	Components []Component
}

// Component is the component name and line upon which it was found
type Component struct {
	Line            int
	Name            string
	HasSerializer   bool
	HasDeserializer bool
}

func main() {
	var vars Vars
	flag.StringVar(&vars.Package, "package", "", "the package where this component will live")
	flag.StringVar(&vars.File, "file", "", "the source file that contains the '//go:generate escgen ...'")
	flag.Parse()

	fmt.Printf("ecsgen v%s: package=\"%s\", file=\"%s\"\n", version, vars.Package, vars.File)

	if len(vars.Package) == 0 || len(vars.File) == 0 {
		fmt.Printf("missing parameter(s)\n")
		os.Exit(-1)
	}

	// we know some of the template vars already
	vars.Backtick = "`"
	vars.Version = version

	{
		// open the source go file to retreive the componant type name
		f, err := os.Open(vars.File)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		vars.Components = getComponents(f)
		if len(vars.Components) == 0 {
			panic(fmt.Errorf("No components found in file %s", vars.File))
		}
	}

	// make generated file name
	vars.GenFile = getOutputFileName(vars.File)

	// now generate the output from the template
	tmpl, err := template.New("component").Parse(templatedOutput)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(vars.GenFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, vars)
	if err != nil {
		panic(err)
	}
}
