package engine

import (
	"container/heap"
	"container/list"
	"fmt"
)

var mdl *model

type event struct {
	time    float64
	toOwner chan int
}

type priorityQueue []*event

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].time < pq[j].time // greater time - less priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*event)
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
	time            float64
	eventChan       chan *event
	queueChan       chan (chan int)
	actorCount      int
	pq              priorityQueue
	bookkeeping     Stats
}

func newModel() *model {
	m := &model{}
	m.blockedInQueues = list.New()
	m.waiting = list.New()
	m.eventChan = make(chan *event)
	m.queueChan = make(chan (chan int))
	m.pq = make(priorityQueue, 0)
	heap.Init(&m.pq)
	return m
}

type ActorInterface interface {
	Run()
	GetGenericActor() *Actor
}

func (m *model) RegisterActor(a ActorInterface) {
	genericActor := a.GetGenericActor()
	genericActor.toModelEvent = m.eventChan
	genericActor.toModelQueue = m.queueChan
	m.actorCount += 1

	go a.Run()
}

func (m *model) getTime() float64 {
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

func (m *model) run(threshold float64) {
	//wait for all actors to start and add an event or block on a queue
	for i := 0; i < m.actorCount; i++ {
		m.waitActor()
	}
	fmt.Printf("All actors started\n")

	//all actors started
	for m.time < threshold {

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
		event := heap.Pop(&m.pq).(*event)
		m.time = event.time
		event.toOwner <- 1

		// wait till process adds event or blocks in queue
		m.waitActor()
	}
	m.bookkeeping.PrintStats()
}

//FIXME remove integers with real requests or sth generic
type QueueInterface interface {
	Enqueue(interface{})
	Dequeue() interface{}
	Len() int
}

type Actor struct {
	toModelEvent chan *event
	toModelQueue chan (chan int)
	inQueue      QueueInterface
	outQueue     QueueInterface
}

func (a *Actor) SetInQueue(q QueueInterface) {
	a.inQueue = q
}

func (a *Actor) SetOutQueue(q QueueInterface) {
	a.outQueue = q
}

func (a *Actor) Wait(d float64) {
	e := &event{time: d + mdl.getTime()}
	ch := make(chan int)
	e.toOwner = ch
	a.toModelEvent <- e
	<-ch // block
}

func (a *Actor) ReadInQueue() interface{} {
	if a.inQueue.Len() > 0 {
		return a.inQueue.Dequeue()
	}
	ch := make(chan int)
	a.toModelQueue <- ch
	<-ch
	return a.ReadInQueue()
}

func (a *Actor) WriteOutQueue(el interface{}) {
	a.outQueue.Enqueue(el)
}

func InitSim() {
	mdl = newModel()
}

func GetTime() float64 {
	return mdl.getTime()
}

func RegisterActor(a ActorInterface) {
	mdl.RegisterActor(a)
}

func Run(threshold float64) {
	mdl.run(threshold)
}

type Stats interface {
	PrintStats()
}

func InitStats(s Stats) {
	mdl.bookkeeping = s
}
