package blocks

import (
	"container/list"
	//	"fmt"

	"github.com/marioskogias/schedsim/engine"
)

type Processor interface {
	engine.ActorInterface
	SetReqDrain(rd RequestDrain) // We might want to specify different drains for different processors or use the same drain for all
	SetCtxCost(cost float64)
}

type RequestDrain interface {
	TerminateReq(r Request)
}

// generic processor: All processors should have it as an embedded field
type genericProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	ctxCost  float64
}

func (p *genericProcessor) GetGenericActor() *engine.Actor {
	return &p.Actor
}

func (p *genericProcessor) SetReqDrain(rd RequestDrain) {
	p.reqDrain = rd
}

func (p *genericProcessor) SetCtxCost(cost float64) {
	p.ctxCost = cost
}

// Run to completion processor
type RTCProcessor struct {
	genericProcessor
	scale float64
}

func (p *RTCProcessor) Run() {
	for {
		//		t1 := engine.GetTime()
		req := p.ReadInQueue().(Request)
		//		t2 := engine.GetTime()
		//		fmt.Printf("%v\n", t2-t1)
		//fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())
		p.Wait(req.ServiceTime + p.ctxCost)
		p.reqDrain.TerminateReq(req)
	}
}

// Time sharing processor
type TSProcessor struct {
	genericProcessor
	quantum float64
}

func NewTSProcessor(quantum float64) *TSProcessor {
	return &TSProcessor{quantum: quantum}
}

func (p *TSProcessor) Run() {
	for {
		req := p.ReadInQueue().(Request)
		//fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())

		if req.ServiceTime <= p.quantum {
			p.Wait(req.ServiceTime + p.ctxCost)
			p.reqDrain.TerminateReq(req)
		} else {
			p.Wait(p.quantum + p.ctxCost)
			req.ServiceTime -= p.quantum
			p.WriteInQueue(req)
		}
	}
}

// Processor sharing processor
type PSProcessor struct {
	genericProcessor
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

// Run to hybrid processor processor
type HybridProcessor struct {
	genericProcessor
	Threshold float64
}

func NewHybridProcessor(quantum float64) *HybridProcessor {
	return &HybridProcessor{Threshold: quantum}
}

func (p *HybridProcessor) Run() {
	for {
		req := p.ReadInQueue().(Request)
		//fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())
		if req.ServiceTime <= p.Threshold {
			p.Wait(req.ServiceTime + p.ctxCost)
			p.reqDrain.TerminateReq(req)
		} else {
			p.Wait(p.Threshold + p.ctxCost)
			req.ServiceTime -= p.Threshold
			p.WriteOutQueue(req)
		}
	}
}

// QoS processor
type QoSProcessor struct {
	genericProcessor
	reqDrains []RequestDrain
}

func (p *QoSProcessor) SetReqDrain(rd RequestDrain) {
	p.reqDrains = append(p.reqDrains, rd)
}

func (p *QoSProcessor) Run() {
	for {
		reqI, _ := p.ReadInQueuesW()
		req := reqI.(Request)
		p.Wait(req.ServiceTime + p.ctxCost)
		p.reqDrains[req.QoS].TerminateReq(req)
	}
}
