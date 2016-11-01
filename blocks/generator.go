package blocks

import (
	"fmt"

	"github.com/marioskogias/schedsim2/engine"
)

type Generator struct {
	engine.Actor
}

func (a *Generator) Run() {
	for {
		fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		a.WriteOutQueue(1)
		a.WriteOutQueue(1)
		a.Wait(5)
	}

}

func (a *Generator) GetGenericActor() *engine.Actor {
	return &a.Actor
}
