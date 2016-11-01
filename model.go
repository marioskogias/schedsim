package main

import (
	"container/heap"
	"container/list"
	"fmt"
)

var mdl *model

type Event struct {
	time    int
	toOwner chan int
}

type priorityQueue []*Event

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].time < pq[j].time // greater time - less priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*Event)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	event := old[n-1]
	*pq = old[0 : n-1]
	return event
}

type model struct {
	blockedInQueues *list.List
	waiting         *list.List
	time            int
	eventChan       chan *Event
	queueChan       chan (chan int)
	actorCount      int
	pq              priorityQueue
}

func newModel() *model {
	m := &model{}
	m.blockedInQueues = list.New()
	m.waiting = list.New()
	m.eventChan = make(chan *Event)
	m.queueChan = make(chan (chan int))
	m.pq = make(priorityQueue, 0)
	heap.Init(&m.pq)
	return m
}

func (m *model) registerActor(a *actor) {
	a.toModelEvent = m.eventChan
	a.toModelQueue = m.queueChan
	m.actorCount += 1
	go a.run()
}

func (m *model) getTime() int {
	return m.time
}

func (m *model) waitActor() {
	select {
	case event := <-m.eventChan:
		heap.Push(&m.pq, event)
	case blocked := <-m.queueChan:
		m.blockedInQueues.PushBack(blocked)
	}
}

func (m *model) run() {
	//wait for all actors to start and add an event or block on a queue
	for i := 0; i < m.actorCount; i++ {
		m.waitActor()
	}
	fmt.Printf("All actors started\n")

	//all actors started
	for m.time < 10 {

		//Check blocked in queues
		if m.blockedInQueues.Len() > 0 {
			l := m.blockedInQueues
			m.blockedInQueues = list.New()

			for e := l.Front(); e != nil; e = e.Next() {
				ch := e.Value.(chan int)
				ch <- 1 // try to unblock
				//wait to block again
				m.waitActor()
			}

		}
		// pick event and wake up process
		event := heap.Pop(&m.pq).(*Event)
		m.time = event.time
		event.toOwner <- 1

		// wait till process adds event or blocks in queue
		m.waitActor()
	}
}

type queue struct {
	l *list.List
}

func newQueue() *queue {
	q := &queue{}
	q.l = list.New()
	return q
}

func (q *queue) enqueue(el int) {
	q.l.PushBack(el)
}

func (q *queue) dequeue() int {
	el := q.l.Front()
	q.l.Remove(el)
	return el.Value.(int)
}

func (q *queue) len() int {
	return q.l.Len()
}

type actor struct {
	toModelEvent chan *Event
	toModelQueue chan (chan int)
	inQueue      *queue
	outQueue     *queue
	name         string
}

func (a *actor) setInQueue(q *queue) {
	a.inQueue = q
}

func (a *actor) setOutQueue(q *queue) {
	a.outQueue = q
}

func (a *actor) wait(d int) {
	e := &Event{time: d + mdl.getTime()}
	ch := make(chan int)
	e.toOwner = ch
	a.toModelEvent <- e
	<-ch // block
}

func (a *actor) readInQueue() int {
	if a.inQueue.len() > 0 {
		return a.inQueue.dequeue()
	}
	ch := make(chan int)
	a.toModelQueue <- ch
	<-ch
	return a.readInQueue()
}

func (a *actor) run() {
	if a.name == "generator" {
		for {
			fmt.Printf("Generator: will add in queue TIME = %v\n", mdl.getTime())
			a.outQueue.enqueue(1)
			a.outQueue.enqueue(1)
			a.wait(5)
		}
	} else {
		for {
			req := a.readInQueue()
			fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req, mdl.getTime())
			a.wait(req)
		}
	}
}

func main() {
	mdl = newModel()

	generator := &actor{name: "generator"}
	processor := &actor{name: "processor"}
	q := newQueue()

	generator.setOutQueue(q)
	processor.setInQueue(q)

	mdl.registerActor(generator)
	mdl.registerActor(processor)
	mdl.run()

}
