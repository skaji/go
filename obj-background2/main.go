package main

import "time"

// Foo is
type Foo struct {
	stop chan<- struct{}
	done <-chan struct{}
}

// NewFoo is
func NewFoo() *Foo {
	f := &Foo{}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		f.background(stop)
	}()
	f.stop = stop
	f.done = done
	return f
}

// Close is
func (f *Foo) Close() {
	close(f.stop)
	<-f.done
}

func (f *Foo) background(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			//
			time.Sleep(time.Second)
		}
	}
}

func main() {
	f := NewFoo()
	f.Close()
}
