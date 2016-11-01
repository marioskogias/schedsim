package blocks

import (
	"fmt"

	"github.com/marioskogias/schedsim2/engine"
)

type Processor struct {
	engine.Actor
}

func (a *Processor) Run() {
	for {
		req := a.ReadInQueue().(Request)
		fmt.Printf("Processor: read from queue val = %v TIME = %v\n", req.ServiceTime, engine.GetTime())
		a.Wait(req.ServiceTime)
	}

}

func (a *Processor) GetGenericActor() *engine.Actor {
	return &a.Actor
}
