package blocks

import (
	"fmt"

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
	// sth to keep the requets
}

func NewPSProcessor() *PSProcessor {
	return &PSProcessor{}
}

func (p *PSProcessor) Run() {
	var d float64
	d = -1
	for {
		intr, reqIntrf := p.ReadInQueueTimeOut(d)
		if intr {
			fmt.Printf("Timeout triggered\n")
			d = -1
		} else {
			fmt.Printf("New request came\n")
			req := reqIntrf.(Request)
			d = req.ServiceTime
			fmt.Printf("The service time is %v\n", d)
		}
	}
}

func (a *PSProcessor) GetGenericActor() *engine.Actor {
	return &a.Actor
}

func (a *PSProcessor) SetReqDrain(rd RequestDrain) {
	a.reqDrain = rd
}
