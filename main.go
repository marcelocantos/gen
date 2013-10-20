package main

import (
	"bitbucket.org/pkg/inflect"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"
)

type Values struct {
	Package   string
	Singular  string
	Plural    string
	Receiver  string
	Loop      string
	Pointer   string
	Generated string
	Command   string
	FileName  string
}

func main() {
	has_args := len(os.Args) > 1
	if !has_args {
		fmt.Print(usage)
		return
	}

	v := getValues()
	t := getTemplate()
	writeFile(t, v)
}

var arg = regexp.MustCompile(`(\*?)(\p{L}+)\.(\p{L}+)`)

func getValues() (v *Values) {
	matches := arg.FindStringSubmatch(os.Args[1])

	if matches == nil {
		log.Fatalln("The first argument must be in the form of package.TypeName")
	}

	ptr := matches[1]
	pkg := matches[2]
	typ := inflect.Singularize(matches[3])
	first, _ := utf8.DecodeRuneInString(typ)
	rcv := strings.ToLower(string(first))

	return &Values{
		Package:   pkg,
		Singular:  typ,
		Plural:    inflect.Pluralize(typ),
		Receiver:  rcv,
		Loop:      "_" + rcv,
		Pointer:   ptr,
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   strings.Join(os.Args, " "),
		FileName:  strings.ToLower(typ) + "_gen.go",
	}
}

func writeFile(t *template.Template, v *Values) {
	f, err := os.Create(v.FileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t.Execute(f, v)
}

const usage = `Usage: [*]package.TypeName

* is recommended but optional, and indicates that generated code should use pointers to the type.

This preference is best for implementing 'expected' and more performant behavior; the non-pointer version will copy structs by value with each function call.
`