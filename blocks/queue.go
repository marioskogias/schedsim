package blocks

import (
	"container/list"
)

type Queue struct {
	l *list.List
}

func NewQueue() *Queue {
	q := &Queue{}
	q.l = list.New()
	return q
}

func (q *Queue) Enqueue(el interface{}) {
	q.l.PushBack(el)
}

func (q *Queue) Dequeue() interface{} {
	el := q.l.Front()
	q.l.Remove(el)
	return el.Value
}

func (q *Queue) Len() int {
	return q.l.Len()
}
