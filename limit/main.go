package main

import (
	"context"
	"sync"
)

type Task struct{}

func run(ctx context.Context, task *Task) error {
	return nil
}

func runLimit1(ctx context.Context, num int, tasks []*Task) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		taskCh := make(chan *Task, len(tasks))
		for _, task := range tasks {
			taskCh <- task
		}
		var wg sync.WaitGroup
		for i := 0; i < num; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case t, ok := <-taskCh:
						if !ok {
							return
						}
						if err := run(ctx, t); err != nil {
							errCh <- err
						}
					}
				}
			}()
		}
		wg.Wait()
	}()

	var lastErr error
	for err := range errCh {
		lastErr = err
	}
	return lastErr
}

func runLimit2(ctx context.Context, num int, tasks []*Task) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		limit := make(chan struct{}, num)
		for i := 0; i < num; i++ {
			limit <- struct{}{}
		}
		var wg sync.WaitGroup
		for _, t := range tasks {
			select {
			case <-ctx.Done():
				break
			case <-limit:
			}
			wg.Add(1)
			go func(t *Task) {
				defer func() {
					wg.Done()
					limit <- struct{}{}
				}()
				if err := run(ctx, t); err != nil {
					errCh <- err
				}
			}(t)
		}
		wg.Done()
	}()

	var lastErr error
	for err := range errCh {
		lastErr = err
	}
	return lastErr
}

func main() {
}
