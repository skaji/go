package main

import (
	"log"
	"testing"
)

var dispatch = map[string]func([]string) int{
	"foo": func(args []string) int {
		log.Println("args -> ", args)
		return 0
	},
	"bar": func(args []string) int {
		return 1
	},
	"baz": func(args []string) int {
		return 0
	},
}

func TestDispatch(t *testing.T) {
	osArgs := []string{"program", "oops"}
	if len(osArgs) < 2 {
		log.Fatal("oops")
	}

	subcmd := osArgs[1]
	args := osArgs[2:]

	if fn, ok := dispatch[subcmd]; ok {
		exit := fn(args)
		log.Println("--> exit", exit)
	} else {
		log.Println("unknown subcmd:", subcmd)
	}

}
