package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jessevdk/go-flags"
	"github.com/lestrrat/go-jshschema"
	"github.com/meta-closure/go-jstr"
	_ "io"
	"io/ioutil"
	"log"
	"os"
)

type options struct {
	YAML    string `short:"y" long:"yaml" description:"the source JSON Schema written by YAML"`
	JSON    string `short:"j" long:"json" description:"the source JSON Schema written by JSON"`
	OutFile string `short:"o" long:"out" description:"output file"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}

	if opts.YAML == "" && opts.JSON == "" {
		log.Println("Not specifies JSON Schema file")
		return 1
	}

	var s map[string]interface{}
	if opts.YAML != "" {
		f, err := ioutil.ReadFile(opts.YAML)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}

		err = yaml.Unmarshal(f, &s)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}

	} else {
		f, err := os.Open(opts.JSON)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
		defer f.Close()

		err = json.NewDecoder(f).Decode(&s)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
	}

	hs := hschema.New()
	err := hs.Extract(s)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	b, err := jstr.Generate(hs)
	if err != nil {
		log.Printf("%s", err)
	}

	out := os.Stdout

	if file := opts.OutFile; file != "" {
		out, err = os.Create(file)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
		defer out.Close()
	}

	fmt.Fprintf(out, "%s", b)
	return 0
}
