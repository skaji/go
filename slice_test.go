package main

import (
	"log"
	"testing"
)

// appending item to nil slices is valid
func TestNilSlice(t *testing.T) {
	var s []string
	log.Println(s == nil) // true
	s = append(s, "foo")
	log.Println(s) // ["foo"]

}

// Nomarly, don't pass slices as "value"
// http://jxck.hatenablog.com/entry/golang-slice-internals2
func TestSliceAsValue(t *testing.T) {
	add := func(s []string, v string) {
		s = append(s, v)
	}

	s := make([]string, 0, 10)
	log.Println(s)
	add(s, "foo")
	log.Println(s)          // []
	log.Println(s[:cap(s)]) // [foo  ]
}
