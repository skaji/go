package queue

import (
	"github.com/skaji/go/queue/queue/base"
	"github.com/skaji/go/queue/task"
)

type Queue struct {
	*base.Queue

	dones map[string]struct{}
}

func New() *Queue {
	dones := make(map[string]struct{})
	satisfied := func(t *task.Task) bool {
		for _, dep := range t.Deps {
			if _, ok := dones[dep]; !ok {
				return false
			}
		}
		return true
	}
	wait := func(tasks []*task.Task) bool {
		for _, t := range tasks {
			if satisfied(t) {
				return false
			}
		}
		return true
	}
	shift := func(tasks []*task.Task) (*task.Task, []*task.Task) {
		for i, t := range tasks {
			if satisfied(t) {
				x := tasks[i]
				xs := append(tasks[:i], tasks[i+1:]...)
				return x, xs
			}
		}
		panic("unexpected")
	}
	return &Queue{dones: dones, Queue: base.NewQueue(wait, shift)}
}

func (q *Queue) Done(t *task.Task) {
	q.Queue.Cond.L.Lock()
	q.dones[t.ID] = struct{}{}
	q.Queue.Cond.L.Unlock()
	q.Queue.Cond.Broadcast()
}
