package blocks

import (
	"container/list"

	"github.com/marioskogias/schedsim/engine"
)

type RequestDrain interface {
	TerminateReq(r Request)
}

// Run to completion processor
type RTCProcessor struct {
	engine.Actor
	reqDrain RequestDrain
}

func (p *RTCProcessor) Run() {
	for {
		req := p.ReadInQueue().(Request)
		//fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())
		p.Wait(req.ServiceTime)
		p.reqDrain.TerminateReq(req)
	}
}

func (p *RTCProcessor) GetGenericActor() *engine.Actor {
	return &p.Actor
}

func (p *RTCProcessor) SetReqDrain(rd RequestDrain) {
	p.reqDrain = rd
}

// Time sharing processor
type TSProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	quantum  float64
}

func NewTSProcessor(quantum float64) *TSProcessor {
	return &TSProcessor{quantum: quantum}
}

func (p *TSProcessor) Run() {
	for {
		req := p.ReadInQueue().(Request)
		//fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())

		if req.ServiceTime <= p.quantum {
			p.Wait(req.ServiceTime)
			p.reqDrain.TerminateReq(req)
		} else {
			p.Wait(p.quantum)
			req.ServiceTime -= p.quantum
			p.WriteInQueue(req)
		}
	}
}

func (a *TSProcessor) GetGenericActor() *engine.Actor {
	return &a.Actor
}

func (a *TSProcessor) SetReqDrain(rd RequestDrain) {
	a.reqDrain = rd
}

// Processor sharing processor
type PSProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	count    int // how many concurrent requests
	reqList  *list.List
	curr     *list.Element
	prevTime float64
}

func NewPSProcessor() *PSProcessor {
	return &PSProcessor{reqList: list.New()}
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

func (p *PSProcessor) updateServiceTimes() {
	currTime := engine.GetTime()
	diff := (currTime - p.prevTime) / float64(p.count)
	//fmt.Printf("Diff = %v\n", diff)
	p.prevTime = currTime
	for e := p.reqList.Front(); e != nil; e = e.Next() {
		req := e.Value.(*Request)
		//fmt.Printf("update: ServiceTime=%v, diff = %v\n", req.ServiceTime, diff)
		req.ServiceTime -= diff
		if e.Value.(*Request).ServiceTime < 0 {
			if e != p.curr {
				panic("updateServiceTime is wrong: negative\n")
			}
		}
	}
}

func (p *PSProcessor) Run() {
	var d float64
	d = -1
	for {
		intr, reqIntrf := p.ReadInQueueTimeOut(d)
		//update times
		p.updateServiceTimes()
		if intr {
			req := p.curr.Value.(*Request)
			p.reqDrain.TerminateReq(*req)
			p.reqList.Remove(p.curr)
			p.count--
		} else {
			p.count++
			req := reqIntrf.(Request)

			reqPtr := &req
			p.reqList.PushBack(reqPtr)
		}
		if p.count > 0 {
			p.curr = p.getMinService()
			d = p.curr.Value.(*Request).ServiceTime * float64(p.count)
		} else {
			d = -1
		}
	}
}

func (a *PSProcessor) GetGenericActor() *engine.Actor {
	return &a.Actor
}

func (a *PSProcessor) SetReqDrain(rd RequestDrain) {
	a.reqDrain = rd
}
