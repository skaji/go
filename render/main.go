package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func render(tmpl string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	var f struct{ Foo string }
	f.Foo = "fff"
	a, _ := render("foo bar {{ .Foo }}", f)
	fmt.Println(a)
}
