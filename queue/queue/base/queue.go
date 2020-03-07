package base

import (
	"sync"

	"github.com/skaji/go/queue/task"
)

type WaitFunc func(tasks []*task.Task) bool
type ShiftFunc func(tasks []*task.Task) (*task.Task, []*task.Task)

type Queue struct {
	C      <-chan *task.Task
	queue  []*task.Task
	Cond   *sync.Cond
	closed bool
	done   <-chan struct{}
}

func NewQueue(wait WaitFunc, shift ShiftFunc) *Queue {
	c := make(chan *task.Task)
	done := make(chan struct{})
	q := &Queue{
		C:      c,
		queue:  []*task.Task{},
		Cond:   sync.NewCond(new(sync.Mutex)),
		closed: false,
		done:   done,
	}
	go func() {
		defer func() {
			close(c)
			close(done)
		}()
		for {
			q.Cond.L.Lock()
			for wait(q.queue) && !q.closed {
				q.Cond.Wait()
			}
			if q.closed {
				q.Cond.L.Unlock()
				return
			}
			var t *task.Task
			t, q.queue = shift(q.queue)
			q.Cond.L.Unlock()
			c <- t
		}
	}()
	return q
}

func (q *Queue) Put(t *task.Task) {
	q.Cond.L.Lock()
	q.queue = append(q.queue, t)
	q.Cond.L.Unlock()
	q.Cond.Broadcast()
}

func (q *Queue) Close() {
	q.Cond.L.Lock()
	q.closed = true
	q.Cond.L.Unlock()
	q.Cond.Broadcast()
	<-q.done
}
