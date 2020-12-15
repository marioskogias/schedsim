package blocks

import (
	"container/heap"
	"container/list"
	"math"
	"math/rand"

	"github.com/epfl-dcsl/schedsim/engine"
)

// Processor Interface describes the main processor functionality used
// in describing a topology
type Processor interface {
	engine.ActorInterface
	SetReqDrain(rd RequestDrain) // We might want to specify different drains for different processors or use the same drain for all
	SetCtxCost(cost float64)
}

// generic processor: All processors should have it as an embedded field
type genericProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	ctxCost  float64
}

func (p *genericProcessor) SetReqDrain(rd RequestDrain) {
	p.reqDrain = rd
}

func (p *genericProcessor) SetCtxCost(cost float64) {
	p.ctxCost = cost
}

// RTCProcessor is a run to completion processor
type RTCProcessor struct {
	genericProcessor
	scale float64
}

// Run is the main processor loop
func (p *RTCProcessor) Run() {
	for {
		req := p.ReadInQueue()
		p.Wait(req.GetServiceTime() + p.ctxCost)
		if monitorReq, ok := req.(*MonitorReq); ok {
			monitorReq.finalLength = p.GetInQueueLen(0)
		}
		p.reqDrain.TerminateReq(req)
	}
}

// TSProcessor is a time sharing processor
type TSProcessor struct {
	genericProcessor
	quantum float64
}

// NewTSProcessor returns a new *TSProcessor
func NewTSProcessor(quantum float64) *TSProcessor {
	return &TSProcessor{quantum: quantum}
}

// Run is the main processor loop
func (p *TSProcessor) Run() {
	for {
		req := p.ReadInQueue()

		if req.GetServiceTime() <= p.quantum {
			p.Wait(req.GetServiceTime() + p.ctxCost)
			p.reqDrain.TerminateReq(req)
		} else {
			p.Wait(p.quantum + p.ctxCost)
			req.SubServiceTime(p.quantum)
			p.WriteInQueue(req)
		}
	}
}

// PSProcessor is a processor sharing processor
type PSProcessor struct {
	genericProcessor
	workerCount int
	count       int // how many concurrent requests
	reqList     *list.List
	curr        *list.Element
	prevTime    float64
}

// NewPSProcessor returns a new *PSProcessor
func NewPSProcessor() *PSProcessor {
	return &PSProcessor{workerCount: 1, reqList: list.New()}
}

// SetWorkerCount sets the number of workers in a processor sharing processor
func (p *PSProcessor) SetWorkerCount(count int) {
	p.workerCount = count
}

func (p *PSProcessor) getMinService() *list.Element {
	minS := p.reqList.Front().Value.(*Request).ServiceTime
	minI := p.reqList.Front()
	for e := p.reqList.Front(); e != nil; e = e.Next() {
		val := e.Value.(*Request).ServiceTime
		if val < minS {
			minS = val
			minI = e
		}
	}
	return minI
}

func (p *PSProcessor) getFactor() float64 {
	if p.workerCount > p.count {
		return 1.0
	}
	return float64(p.workerCount) / float64(p.count)
}

func (p *PSProcessor) updateServiceTimes() {
	currTime := engine.GetTime()
	diff := (currTime - p.prevTime) * p.getFactor()
	p.prevTime = currTime
	for e := p.reqList.Front(); e != nil; e = e.Next() {
		req := e.Value.(engine.ReqInterface)
		req.SubServiceTime(diff)
	}
}

// Run is the main processor loop
func (p *PSProcessor) Run() {
	var d float64
	d = -1
	for {
		intr, newReq := p.WaitInterruptible(d)
		//update times
		p.updateServiceTimes()
		if intr {
			req := p.curr.Value.(engine.ReqInterface)
			p.reqDrain.TerminateReq(req)
			p.reqList.Remove(p.curr)
			p.count--
		} else {
			p.count++
			p.reqList.PushBack(newReq)
		}
		if p.count > 0 {
			p.curr = p.getMinService()
			d = p.curr.Value.(engine.ReqInterface).GetServiceTime() / p.getFactor()
		} else {
			d = -1
		}
	}
}

type BoundedProcessor struct {
	genericProcessor
	bufSize int
}

func NewBoundedProcessor(bufSize int) *BoundedProcessor {
	return &BoundedProcessor{bufSize: bufSize}
}

// Run is the main processor loop
func (p *BoundedProcessor) Run() {
	var factor float64
	for {
		req := p.ReadInQueue()

		if colorReq, ok := req.(*ColoredReq); ok {
			if colorReq.color == 1 {
				factor = 2
			} else {
				factor = 1
			}
		}
		p.Wait(factor * req.GetServiceTime())
		len := p.GetOutQueueLen(0)
		if len < p.bufSize {
			p.WriteOutQueue(req)
		} else {
			p.reqDrain.TerminateReq(req)
		}
	}
}

type BoundedProcessor2 struct {
	genericProcessor
}

// Run is the main processor loop
func (p *BoundedProcessor2) Run() {
	var factor float64
	for {
		req := p.ReadInQueue()

		if colorReq, ok := req.(*ColoredReq); ok {
			if colorReq.color == 0 {
				factor = 2
			} else {
				factor = 1
			}
		}
		p.Wait(factor * req.GetServiceTime())
		p.reqDrain.TerminateReq(req)
	}
}

// A verona processor is three things:
// time sharing to emulate a task with many behaviours
// stealing when idle
// fair with the token stealing mechanism (that's hard)
type VeronaProcessor struct {
	genericProcessor
	quantum   float64
	nextSteal int
}

func NewVeronaProcessor(q float64) *VeronaProcessor {
	return &VeronaProcessor{quantum: q}
}

func (p *VeronaProcessor) Run() {
	var r engine.ReqInterface
	var gotReq bool

	for {
		if p.nextSteal == 0 {
			// Avoid queue 0 which is the local queue
			base := rand.Intn(p.GetInQueueCount())
			for i := 0; i < p.GetInQueueCount(); i += 1 {
				idx := (base + i) % p.GetInQueueCount()
				if idx == 0 {
					continue
				}
				l := p.GetInQueueLen(i)
				if l > 0 {
					r = p.ReadInQueueI(i)
					gotReq = true
					break
				}
			}
			p.nextSteal = int(math.Max(1.0, float64(p.GetInQueueLen(0))))
		}
		if !gotReq {
			localCount := p.GetInQueueLen(0)
			if localCount > 0 {
				r = p.ReadInQueueI(0)
			} else {
				r, _ = p.ReadInQueuesRandLocalPr()
			}
		}

		gotReq = false
		// Serve the request
		if r.GetServiceTime() <= p.quantum {
			p.Wait(r.GetServiceTime())
			p.reqDrain.TerminateReq(r)
		} else {
			p.Wait(p.quantum)
			r.SubServiceTime(p.quantum)
			p.WriteInQueue(r)
		}
		p.nextSteal -= 1
	}
}

type VeronaProcessor2 struct {
	genericProcessor
	bCount    int
	nextSteal int
	fair      bool
}

func NewVeronaProcessor2(bCount int, fair bool) *VeronaProcessor2 {
	return &VeronaProcessor2{bCount: bCount, fair: fair}
}

func (p *VeronaProcessor2) Run() {
	var r engine.ReqInterface
	var gotReq bool

	for {
		if p.nextSteal == 0 {
			// Avoid queue 0 which is the local queue
			base := rand.Intn(p.GetInQueueCount())
			for i := 0; i < p.GetInQueueCount(); i += 1 {
				idx := (base + i) % p.GetInQueueCount()
				if idx == 0 {
					continue
				}
				l := p.GetInQueueLen(i)
				if l > 0 {
					r = p.ReadInQueueI(i)
					gotReq = true
					break
				}
			}
			p.nextSteal = int(math.Max(1.0, float64(p.GetInQueueLen(0))))
		}
		if !gotReq {
			localCount := p.GetInQueueLen(0)
			if localCount > 0 {
				r = p.ReadInQueueI(0)
			} else {
				r, _ = p.ReadInQueuesRandLocalPr()
			}
		}

		gotReq = false

		c := r.(*Cown)
		for i := 0; i < p.bCount; i++ {
			if c.queue.Len() == 0 {
				c.isSchedulled = false
				break
			}
			// Serve requests
			el := c.queue.Remove(c.queue.Front())
			req := el.(engine.ReqInterface)
			p.Wait(req.GetServiceTime())
			p.reqDrain.TerminateReq(req)
		}
		if c.isSchedulled {
			p.WriteInQueue(c)
		}
		if p.fair == true {
			p.nextSteal -= 1
		} else {
			p.nextSteal = 100 // avoid stealing
		}
	}
}

// SRPT Processor keeping requests in a heap

type ReqHeap []engine.ReqInterface

func (h ReqHeap) Len() int           { return len(h) }
func (h ReqHeap) Less(i, j int) bool { return h[i].GetServiceTime() < h[j].GetServiceTime() }
func (h ReqHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ReqHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(engine.ReqInterface))
}

func (h *ReqHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type SRPTProcessor struct {
	genericProcessor
	workerCount int
	reqHeap     ReqHeap
	activeList  []engine.ReqInterface
	curr        engine.ReqInterface
	currIdx     int
	prevTime    float64
}

func NewSRPTProcessor(workerCount int) *SRPTProcessor {
	p := &SRPTProcessor{}
	p.workerCount = workerCount
	p.activeList = make([]engine.ReqInterface, workerCount)
	heap.Init(&p.reqHeap)

	return p
}

func (p *SRPTProcessor) updateTimes() {
	currTime := engine.GetTime()
	diff := currTime - p.prevTime
	p.prevTime = currTime
	for _, r := range p.activeList {
		if r != nil {
			r.SubServiceTime(diff)
		}
	}
}

func (p *SRPTProcessor) Run() {
	var d float64
	d = -1
	for {
		intr, newReq := p.WaitInterruptible(d)
		p.updateTimes()
		if intr {
			// If a request finished substitute it from the heap
			p.reqDrain.TerminateReq(p.curr)
			if p.reqHeap.Len() == 0 {
				p.activeList[p.currIdx] = nil
			} else {
				el := heap.Pop(&p.reqHeap)
				p.activeList[p.currIdx] = el.(engine.ReqInterface)
			}
		} else {
			// if new request check if should substitute any of the ones currently running
			// else push on the heap (-1) or put it in the active list if slot is found
			subIdx := -1
			val := 0.0
			for i, r := range p.activeList {
				// put it in the list if a slot is found
				if r == nil {
					p.activeList[i] = newReq
					subIdx = -2
					break
				}
				if r.GetServiceTime() > newReq.GetServiceTime() && val < r.GetServiceTime() {
					val = r.GetServiceTime()
					subIdx = i
				}
			}
			if subIdx > -1 {
				heap.Push(&p.reqHeap, p.activeList[subIdx])
				p.activeList[subIdx] = newReq
			} else if subIdx == -1 {
				heap.Push(&p.reqHeap, newReq)
			}
		}
		// Find the next request to wait for from active list
		idx := -1
		val := 100000.0
		for i, r := range p.activeList {
			if r == nil {
				continue
			}
			if r.GetServiceTime() < val {
				idx = i
				val = r.GetServiceTime()
			}
		}
		if idx > -1 {
			d = p.activeList[idx].GetServiceTime()
			p.curr = p.activeList[idx]
			p.currIdx = idx
		} else {
			d = -1
		}
	}
}
