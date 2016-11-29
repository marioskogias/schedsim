package blocks

import (
	"container/heap"
	"container/list"
	//"sort"
	"fmt"
	//"github.com/marioskogias/schedsim/engine"
)

var count = 0

// Simple FIFO queue
type Queue struct {
	l  *list.List
	id int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.l = list.New()
	q.id = count
	count++
	return q
}

func (q *Queue) Enqueue(el interface{}) {
	//fmt.Printf("time: %v, queue: %v, len: %v\n", engine.GetTime(), q.id, q.Len())
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
	GetServiceTime() float64
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

func (pq *PQueue) Enqueue(el interface{}) {
	fmt.Printf("%v\t", pq.Len())
	pq.PrintQueue()
	fmt.Printf("\n")
	heap.Push(&pq.pq, el)
}

func (pq *PQueue) Dequeue() interface{} {
	return heap.Pop(&pq.pq)
}

func (pq *PQueue) Len() int {
	return pq.pq.Len()
}

func (pq *PQueue) PrintQueue() {
	for _, v := range pq.pq {
		fmt.Printf("%v\t", v.GetServiceTime())
	}
}
