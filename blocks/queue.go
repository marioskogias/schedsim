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

func (q *Queue) Enqueue(el int) {
	q.l.PushBack(el)
}

func (q *Queue) Dequeue() int {
	el := q.l.Front()
	q.l.Remove(el)
	return el.Value.(int)
}

func (q *Queue) Len() int {
	return q.l.Len()
}
