package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		cancel()
	}()
	if err := runContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runNG(stop <-chan struct{}) error {
	svr := server(80)

	done := make(chan error)
	go func() {
		<-stop
		done <- svr.Shutdown(context.Background())
	}()

	err1 := svr.ListenAndServe()
	err2 := <-done
	if err1 != nil && err1 != http.ErrServerClosed {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func run(stop <-chan struct{}) error {
	svr := server(80)
	done := make(chan error)
	myStop := make(chan struct{})
	go func() {
		select {
		case <-stop:
		case <-myStop:
		}
		done <- svr.Shutdown(context.Background())
	}()
	err1 := svr.ListenAndServe()
	close(myStop)
	err2 := <-done
	if err1 != nil && err1 != http.ErrServerClosed {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func runContext(ctx context.Context) error {
	svr := server(80)
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan error)
	go func() {
		<-ctx.Done()
		done <- svr.Shutdown(context.Background())
	}()
	err1 := svr.ListenAndServe()
	cancel()
	err2 := <-done
	if err1 != nil && err1 != http.ErrServerClosed {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
