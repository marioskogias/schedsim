package engine

import (
	"container/heap"
	"container/list"
	"math/rand"
)

var mdl *model
var Weight float32

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
	active       bool
}

type model struct {
	blockedInQueues *list.List
	waiting         *list.List
	time            float64
	eventChan       chan *event
	queueChan       chan *blockEvent
	actorCount      int
	pq              priorityQueue
	bookkeeping     []Stats
}

func newModel() *model {
	m := &model{}
	m.blockedInQueues = list.New()
	m.waiting = list.New()
	m.eventChan = make(chan *event)
	m.queueChan = make(chan *blockEvent)
	m.pq = make(priorityQueue, 0)
	heap.Init(&m.pq)
	return m
}

type ActorInterface interface {
	Run()
	GetGenericActor() *Actor
	AddInQueue(q QueueInterface)
	AddOutQueue(q QueueInterface)
	GetOutQueueLengths() []int
}

func (m *model) registerActor(a ActorInterface) {
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
	case blocked := <-m.queueChan: // Actor did ReadInqueue: Blocked in queue
		if blocked.timeOutEvent != nil {
			heap.Push(&m.pq, blocked.timeOutEvent)
		}
		m.blockedInQueues.PushBack(blocked)
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
				be := e.Value.(*blockEvent)
				if be.active {
					be.wakeUpCh <- 1 // try to unblock
					//wait to block again
					m.waitActor()
				}
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
	for _, s := range m.bookkeeping {
		s.PrintStats()
	}
}

type QueueInterface interface {
	Enqueue(interface{})
	Dequeue() interface{}
	Len() int
}

type Actor struct {
	toModelEvent chan *event
	toModelQueue chan *blockEvent
	inQueues     []QueueInterface
	outQueues    []QueueInterface
}

// In and out queues should be added in decreasing priority
func (a *Actor) AddInQueue(q QueueInterface) {
	a.inQueues = append(a.inQueues, q)
}

func (a *Actor) AddOutQueue(q QueueInterface) {
	a.outQueues = append(a.outQueues, q)
}

func (a *Actor) GetInQueueLen(idx int) int {
	return a.inQueues[idx].Len()
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
	if a.inQueues[0].Len() > 0 {
		return false, a.inQueues[0].Dequeue()
	}

	// Negative timeout - no timeout
	if d < 0 {
		return false, a.ReadInQueue()
	}
	timeoutTime := d + mdl.getTime()
	e := &event{time: timeoutTime, active: true}
	ch := make(chan int)
	e.toOwner = ch
	bEvent := &blockEvent{timeOutEvent: e, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	for { // this is because the run time tries to run the actors on every iteration
		<-ch
		if a.inQueues[0].Len() > 0 {
			e.active = false
			return false, a.inQueues[0].Dequeue()
		}
		if mdl.getTime() == timeoutTime {
			bEvent.active = false
			return true, nil
		}
		bEvent.timeOutEvent = nil
		a.toModelQueue <- bEvent
	}
}

func (a *Actor) ReadInQueue() interface{} {
	if a.inQueues[0].Len() > 0 {
		return a.inQueues[0].Dequeue()
	}
	ch := make(chan int)
	bEvent := &blockEvent{timeOutEvent: nil, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueue()
}

// This function tries to read from all the queues in descending priority
// and blocks only if all the queues are empty. In returns the element of the
// first queue found non-empty
func (a *Actor) ReadInQueues() (interface{}, int) {
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			return q.Dequeue(), i
		}
	}
	ch := make(chan int)
	bEvent := &blockEvent{timeOutEvent: nil, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueues()
}

type queueIdx struct {
	idx int
	q   QueueInterface
}

func (a *Actor) ReadInQueuesRand() (interface{}, int) {
	var available []queueIdx
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			available = append(available, queueIdx{i, q})
		}
	}
	if len(available) > 0 {
		q := available[rand.Intn(len(available))]
		return q.q.Dequeue(), q.idx
	}
	ch := make(chan int)
	bEvent := &blockEvent{timeOutEvent: nil, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueues()
}

func (a *Actor) ReadInQueuesRandLocalPr() (interface{}, int) {
	if a.inQueues[0].Len() > 0 {
		return a.inQueues[0].Dequeue(), 0
	}
	var available []queueIdx
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			available = append(available, queueIdx{i, q})
		}
	}
	if len(available) > 0 {
		q := available[rand.Intn(len(available))]
		return q.q.Dequeue(), q.idx
	}
	ch := make(chan int)
	bEvent := &blockEvent{timeOutEvent: nil, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueues()
}

// WRR approximation
func (a *Actor) ReadInQueuesW() (interface{}, int) {

	if len(a.inQueues) >= 2 {
		if rand.Float32() >= Weight {
			if a.inQueues[0].Len() > 0 {
				return a.inQueues[0].Dequeue(), 0
			} else {
				if a.inQueues[1].Len() > 0 {
					return a.inQueues[1].Dequeue(), 1
				}
			}
		} else {
			if a.inQueues[1].Len() > 0 {
				return a.inQueues[1].Dequeue(), 1
			} else {
				if a.inQueues[0].Len() > 0 {
					return a.inQueues[0].Dequeue(), 0
				}
			}
		}
		if len(a.inQueues) == 3 {
			if a.inQueues[2].Len() > 0 {
				return a.inQueues[2].Dequeue(), 2

			}
		}
	} else {
		for i, q := range a.inQueues {
			if q.Len() > 0 {
				return q.Dequeue(), i
			}
		}
	}

	ch := make(chan int)
	bEvent := &blockEvent{timeOutEvent: nil, wakeUpCh: ch, active: true}
	a.toModelQueue <- bEvent
	<-ch
	return a.ReadInQueues()
}

func (a *Actor) WriteOutQueue(el interface{}) {
	a.outQueues[0].Enqueue(el)
}

func (a *Actor) WriteInQueue(el interface{}) {
	a.inQueues[0].Enqueue(el)
}

func (a *Actor) WriteOutQueueI(el interface{}, i int) {
	a.outQueues[i].Enqueue(el)
}

func (a *Actor) WriteInQueueI(el interface{}, i int) {
	a.inQueues[i].Enqueue(el)
}

func (a *Actor) GetOutQueueLengths() []int {
	res := make([]int, len(a.outQueues))
	for i, q := range a.outQueues {
		res[i] = q.Len()
	}
	return res
}

func (a *Actor) OutQueueCount() int {
	return len(a.outQueues)
}

func InitSim() {
	mdl = newModel()
}

func GetTime() float64 {
	return mdl.getTime()
}

func RegisterActor(a ActorInterface) {
	mdl.registerActor(a)
}

func Run(threshold float64) {
	mdl.run(threshold)
}

type Stats interface {
	PrintStats()
}

func InitStats(s Stats) {
	mdl.bookkeeping = append(mdl.bookkeeping, s)
}
