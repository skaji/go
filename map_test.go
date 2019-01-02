package main

import (
	"log"
	"sort"
	"testing"
)

func TestNilMap(t *testing.T) {
	defer func() { log.Println("recover", recover()) }()
	var m map[string]string
	log.Println("--->", m["foo"]) // not panic
	m["foo"] = "bar"              // panic: assignment to entry in nil map
}

func TestBasicMap(t *testing.T) {
	m := make(map[string]string, 10)
	m["a"] = "bar"
	m["b"] = "baz"
	m["c"] = "baz"
	log.Println(len(m)) // 3

	// sort keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		log.Println(k, "->", m[k])
	}
}
