package jstr

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lestrrat/go-jshschema"
	"github.com/lestrrat/go-jsschema"
	"go/format"
	"strings"
)

var TypeCorrespond = map[string]string{
	"integer": "MaybeInt",
	"string":  "MaybeString",
	"boolean": "MaybeBool",
	"number":  "MaybeFloat",
}

func Generate(scm *hschema.HyperSchema) ([]byte, error) {

	buf := Init()
	err := _generate(scm.Schema, &buf)
	if err != nil {
		return nil, err
	}
	fbuf, err := format.Source(buf.Bytes())
	return fbuf, nil
}

func _generate(scm *schema.Schema, buf *bytes.Buffer) error {

	for i, j := range scm.Definitions {
		if j.Type == nil {
			err := errors.New("Type unset")
			return err
		}
		switch j.Type[0].String() {
		// example "type Foo struct{
		// 	Foo Bar `json:"foo,omitempty"`
		// }
		case "object":
			{
				fmt.Fprintf(buf, "\ntype %s struct {",
					Snake2Caml(i))

				if j.Properties == nil {
					err := errors.New("Properties unset")
					return err
				}
				for column, ref := range j.Properties {
					columnType, err := ChildSchemaType(ref, scm)
					if err != nil {
						return err
					}
					fmt.Fprintf(buf, "\n %s %s`json:\"%s,omitempty\"`",
						Snake2Caml(column),
						columnType,
						column)
				}
				buf.WriteString("\n}")
			}
		// example "type Foo []Bar"
		case "array":
			{
				refSchemas := j.Items.Schemas
				if refSchemas == nil {
					err := errors.New("Item is unset")
					return err
				}

				columnType, err := ChildSchemaType(refSchemas[0], scm)
				if err != nil {
					return err
				}
				fmt.Fprintf(buf, "\n type %s []%s",
					Snake2Caml(i),
					columnType)
			}
		}
	}
	return nil
}

func ChildSchemaType(ref, root *schema.Schema) (string, error) {
	var str string
	refSchema := root.Definitions[Ref2Name(ref.Reference)]

	switch refSchema.Type[0].String() {
	case "object":
		{
			str = Snake2Caml(refSchema.Title)
		}

	case "array":
		{
			if ref := refSchema.Items.Schemas; ref == nil {
				err := errors.New("Item unset")
				return str, err
			}

			if str, err := ChildSchemaType(ref, root); err != nil {
				return str, err
			}
			str = "[]" + str
		}
	case "string", "integer", "boolean", "number":
		{
			str = "jsval." + TypeCorrespond[refSchema.Type[0].String()]
		}
	default:
		{
			return str, errors.New("Type error")
		}
	}
	return str, nil
}

func Snake2Caml(name string) string {
	var t string
	for _, s := range strings.Split(name, "_") {
		t += strings.Title(s)
	}
	return t
}

func Ref2Name(ref string) string {
	for _, s := range strings.Split(ref, "/") {
		if s != "definitions" && s != "#" {
			return s
		}
	}
	return ""
}

func Init() bytes.Buffer {
	var buf bytes.Buffer
	init := `
	package model
	import (
		"github.com/lestrrat/go-jsval"
	)
	
	`
	buf.WriteString(init)
	return buf
}
