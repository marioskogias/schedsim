package blocks

import (
	"container/heap"
	"container/list"
)

// Simple FIFO queue
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

// PriorityQueue
type Comparable interface {
	GetCmpVal() float64
}

type pQueue []Comparable

func (pq pQueue) Len() int { return len(pq) }

func (pq pQueue) Less(i, j int) bool {
	return pq[i].GetCmpVal() < pq[j].GetCmpVal() // greater time - less priority
}

func (pq pQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pQueue) Push(x interface{}) {
	item := x.(Comparable)
	*pq = append(*pq, item)
}

func (pq *pQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type PQueue struct {
	pq pQueue
}

func NewPQueue() *PQueue {
	q := &PQueue{}
	q.pq = make(pQueue, 0)
	heap.Init(&q.pq)

	return q
}
