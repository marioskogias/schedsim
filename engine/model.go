package engine

import (
	"container/heap"
	"container/list"
)

var mdl *model

type event struct {
	time    float64
	active  bool
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

type blockEvent struct {
	wakeUpCh     chan int
	timeOutEvent *event // if nil no timeout
}

type model struct {
	blockedInQueues *list.List
	waiting         *list.List
	time            float64
	eventChan       chan *event
	queueChan       chan blockEvent
	actorCount      int
	pq              priorityQueue
	bookkeeping     Stats
}

func newModel() *model {
	m := &model{}
	m.blockedInQueues = list.New()
	m.waiting = list.New()
	m.eventChan = make(chan *event)
	m.queueChan = make(chan blockEvent)
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
	case event := <-m.eventChan: // Actor did Wait: new event
		heap.Push(&m.pq, event)
		//FIXME: add timeouts
	case blocked := <-m.queueChan: // Actor did ReadInqueue: Blocked in queue
		if blocked.timeOutEvent != nil {
			heap.Push(&m.pq, blocked.timeOutEvent)
		}
		m.blockedInQueues.PushBack(blocked.wakeUpCh)
	}
}

func (m *model) run(threshold float64) {
	//wait for all actors to start and add an event or block on a queue
	for i := 0; i < m.actorCount; i++ {
		m.waitActor()
	}

	//all actors started
	for m.time < threshold {

		//Check blocked in queues
		if m.blockedInQueues.Len() > 0 {
			l := m.blockedInQueues
			m.blockedInQueues = list.New()

			for e := l.Front(); e != nil; e = e.Next() {
				// FIXME
				ch := e.Value.(chan int)
				ch <- 1 // try to unblock
				//wait to block again
				m.waitActor()
			}
		}
		// pick event and wake up process
		e := heap.Pop(&m.pq).(*event)
		for !e.active {
			e = heap.Pop(&m.pq).(*event)
		}
		m.time = e.time
		e.toOwner <- 1

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
	toModelQueue chan blockEvent
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
	e := &event{time: d + mdl.getTime(), active: true}
	ch := make(chan int)
	e.toOwner = ch
	a.toModelEvent <- e
	<-ch // block
}

// This is not tested. Do we need it?
func (a *Actor) WaitInterruptible(d float64, intr <-chan int) {
	e := &event{time: d + mdl.getTime(), active: true}
	ch := make(chan int)
	e.toOwner = ch
	a.toModelEvent <- e
	select {
	case <-ch:
		return
	case <-intr:
		// Deactivate the event
		e.active = false
	}
}

func (a *Actor) ReadInQueueTimeOut(d float64) (bool, interface{}) {
	if a.inQueue.Len() > 0 {
		return false, a.inQueue.Dequeue()
	}

	timeoutTime := d + mdl.getTime()
	e := &event{time: timeoutTime, active: true}
	ch := make(chan int)
	bEvent := blockEvent{timeOutEvent: e, wakeUpCh: ch}
	a.toModelQueue <- bEvent
	for {
		<-ch
		if a.inQueue.Len() > 0 {
			e.active = false
			return false, a.inQueue.Dequeue()
		}
		if mdl.getTime() == timeoutTime {
			return true, nil
		}
		bEvent := blockEvent{timeOutEvent: nil, wakeUpCh: ch}
		a.toModelQueue <- bEvent
	}
}

func (a *Actor) ReadInQueue() interface{} {
	if a.inQueue.Len() > 0 {
		return a.inQueue.Dequeue()
	}
	ch := make(chan int)
	bEvent := blockEvent{timeOutEvent: nil, wakeUpCh: ch}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueue()
}

func (a *Actor) WriteOutQueue(el interface{}) {
	a.outQueue.Enqueue(el)
}

func (a *Actor) WriteInQueue(el interface{}) {
	a.inQueue.Enqueue(el)
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
