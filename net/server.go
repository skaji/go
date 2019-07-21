package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

// Listener is
type Listener struct {
	net.Listener
}

// AcceptContext is
func (l Listener) AcceptContext(ctx context.Context) (Conn, error) {
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			l.Close()
		case <-stop:
		}
		close(done)
	}()
	c, err := l.Accept()
	close(stop)
	<-done
	return Conn{c}, err
}

// Conn is
type Conn struct {
	net.Conn
}

var zeroTime = time.Time{}
var aLongTimeAgo = time.Unix(1, 0)

// ReadContext is
func (c Conn) ReadContext(ctx context.Context, b []byte) (int, error) {
	c.Conn.SetReadDeadline(zeroTime)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.Conn.SetReadDeadline(aLongTimeAgo)
		case <-stop:
		}
		close(done)
	}()
	n, err := c.Conn.Read(b)
	close(stop)
	<-done
	c.Conn.SetReadDeadline(zeroTime)
	return n, err
}

// WriteContext is
func (c Conn) WriteContext(ctx context.Context, b []byte) (int, error) {
	c.Conn.SetWriteDeadline(zeroTime)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.Conn.SetWriteDeadline(aLongTimeAgo)
		case <-stop:
		}
		close(done)
	}()
	n, err := c.Conn.Write(b)
	close(stop)
	<-done
	c.Conn.SetWriteDeadline(zeroTime)
	return n, err
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		fmt.Println(conn, err)
		if err != nil {
			log.Print(err)
			continue
		}
		go func() {
			defer conn.Close()

			conn.SetReadDeadline(time.Time{})
			go func() {
				time.Sleep(2 * time.Second)
				conn.SetReadDeadline(time.Unix(1, 0))
			}()
			b := make([]byte, 10)
			n, err := conn.Read(b)
			fmt.Println(n, err)
		}()
	}
}
