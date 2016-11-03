package blocks

import (
	//"fmt"

	"github.com/marioskogias/schedsim/engine"
)

type RequestDrain interface {
	TerminateReq(r Request)
}

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

type SharedProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	quantum  float64
}

func NewSharedProcessor(quantum float64) *SharedProcessor {
	return &SharedProcessor{quantum: quantum}
}

func (p *SharedProcessor) Run() {
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

func (a *SharedProcessor) GetGenericActor() *engine.Actor {
	return &a.Actor
}

func (a *SharedProcessor) SetReqDrain(rd RequestDrain) {
	a.reqDrain = rd
}
