package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- "hello"
			time.Sleep(time.Second)
		}
		close(ch1)
	}()
	go func() {
		for i := 0; i < 10; i++ {
			ch2 <- true
			time.Sleep(time.Second)
		}
		close(ch2)
	}()

	for {
		if ch1 == nil && ch2 == nil {
			return
		}
		select {
		case s, ok := <-ch1:
			if ok {
				fmt.Println("got", s)
			} else {
				ch1 = nil
			}
		case b, ok := <-ch2:
			if ok {
				fmt.Println("got", b)
			} else {
				ch2 = nil
			}
		}
	}
}
