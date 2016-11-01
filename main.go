package main

import (
	"github.com/marioskogias/schedsim2/blocks"
	"github.com/marioskogias/schedsim2/engine"
)

func main() {

	engine.InitSim()

	//Add a deterministic generator
	generator := &blocks.Generator{}

	//Add a run to completion processor
	processor := &blocks.Processor{}

	//Add a fifo queue
	q := blocks.NewQueue()

	//Init the statistics
	stats := &blocks.BookKeeper{}
	engine.InitStats(stats)

	//Create the topology
	generator.SetOutQueue(q)
	processor.SetInQueue(q)
	processor.SetReqDrain(stats)

	//Register actors
	engine.RegisterActor(generator)
	engine.RegisterActor(processor)

	//Run till 100000 time units
	engine.Run(10000)
}
