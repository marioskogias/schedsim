package main

import (
	"github.com/marioskogias/schedsim2/blocks"
	"github.com/marioskogias/schedsim2/engine"
)

func main() {

	engine.InitSim()

	generator := &blocks.Generator{}
	processor := &blocks.Processor{}
	q := blocks.NewQueue()

	generator.SetOutQueue(q)
	processor.SetInQueue(q)

	engine.RegisterActor(generator)
	engine.RegisterActor(processor)
	engine.Run(10000)
}
