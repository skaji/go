package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	svr := server(8080)
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

func server(port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := []byte("OK\n")
		w.Header().Set("Connection", "close")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.Header().Set("Content-Type", "text/plain")
		w.Write(body)
	})
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
