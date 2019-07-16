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
	f.stop = stop
	f.done = f.background(stop)
	return f
}

// Close is
func (f *Foo) Close() {
	close(f.stop)
	<-f.done
}

func (f *Foo) background(stop <-chan struct{}) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				//
				time.Sleep(time.Second)
			}
		}
	}()
	return done
}

func main() {
	f := NewFoo()
	f.Close()
}
