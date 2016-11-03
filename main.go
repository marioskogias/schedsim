package main

import (
	"github.com/marioskogias/schedsim2/blocks"
	"github.com/marioskogias/schedsim2/engine"
)

func main() {

	engine.InitSim()

	//Add a deterministic generator
	//generator := blocks.NewDDGenerator(2, 1)

	//Add an MD generator
	//generator := blocks.NewMDGenerator(0.5, 1)

	//Add an MM generator
	generator := blocks.NewMMGenerator(0.005, 0.02) // 50usec sevice time, lambda 0.005

	//Add a run to completion processor
	processor := &blocks.RTCProcessor{}

	//Add a shared processor
	//processor := blocks.NewSharedProcessor(0.5)

	//Add a fifo queue
	q := blocks.NewQueue()

	//Init the statistics
	stats := blocks.NewBookKeeper()
	engine.InitStats(stats)

	//Create the topology
	generator.SetOutQueue(q)
	processor.SetInQueue(q)
	processor.SetReqDrain(stats)

	//Register actors
	engine.RegisterActor(generator)
	engine.RegisterActor(processor)

	//Run till 100000 time units
	engine.Run(10000000)
}
