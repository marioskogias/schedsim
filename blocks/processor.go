package blocks

import (
	"fmt"

	"github.com/marioskogias/schedsim2/engine"
)

type RequestDrain interface {
	TerminateReq(r Request)
}

type Processor struct {
	engine.Actor
	reqDrain RequestDrain
}

func (a *Processor) Run() {
	for {
		req := a.ReadInQueue().(Request)
		fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())
		a.Wait(req.ServiceTime)
		a.reqDrain.TerminateReq(req)
	}
}

func (a *Processor) GetGenericActor() *engine.Actor {
	return &a.Actor
}

func (a *Processor) SetReqDrain(rd RequestDrain) {
	a.reqDrain = rd
}
