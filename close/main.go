package main

import (
	"fmt"
	"time"
)

type foo struct {
	done <-chan struct{}
	stop chan<- struct{}
}

func newFoo() *foo {
	f := &foo{}
	stop := make(chan struct{})
	f.stop = stop
	f.done = f.background(stop)
	return f
}

func (f *foo) close() {
	close(f.stop)
	<-f.done
}

func (f *foo) background(stop <-chan struct{}) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		t := time.NewTicker(time.Second)
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				fmt.Println("alive")
			}
		}
	}()
	return done
}

func main() {
	fmt.Println("create foo")
	f := newFoo()
	fmt.Println("done create foo")
	time.Sleep(5 * time.Second)
	fmt.Println("close foo")
	f.close()
	fmt.Println("done close foo")
}
