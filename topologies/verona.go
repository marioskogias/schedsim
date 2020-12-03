package topologies

import (
	"fmt"

	"github.com/epfl-dcsl/schedsim/blocks"
	"github.com/epfl-dcsl/schedsim/engine"
)

// MultiQueue describes a single-generator-multi-processor topology where every
// processor has its own incoming queue
func Verona(lambda, mu, duration float64, genType int, quantum float64) {

	engine.InitSim()

	//Init the statistics
	//stats := blocks.NewBookKeeper()
	stats := &blocks.AllKeeper{}
	stats.SetName("Main Stats")
	engine.InitStats(stats)

	// Add generator
	var g blocks.Generator
	if genType == 0 {
		g = blocks.NewMMRandGenerator(lambda, mu)
	} else if genType == 1 {
		g = blocks.NewMDRandGenerator(lambda, 1/mu)
	} else if genType == 2 {
		g = blocks.NewMBRandGenerator(lambda, 1, 10*(1/mu-0.9), 0.9)
	} else if genType == 3 {
		g = blocks.NewMBRandGenerator(lambda, 1, 1000*(1/mu-0.999), 0.999)
	}

	g.SetCreator(&blocks.SimpleReqCreator{})

	// Create queues and add them to the generator
	queues := make([]engine.QueueInterface, cores)
	for i := range queues {
		queues[i] = blocks.NewQueue()
		g.AddOutQueue(queues[i])
	}

	// Create processors and assign local queue
	processors := make([]blocks.Processor, cores)
	for i := 0; i < cores; i++ {
		processors[i] = blocks.NewVeronaProcessor(quantum)
		processors[i].AddInQueue(queues[i])
		for j, q := range queues {
			if j == i {
				continue
			} else {
				processors[i].AddInQueue(q)
			}
		}
	}

	// Add the stats and register processors
	for _, p := range processors {
		p.SetReqDrain(stats)
		engine.RegisterActor(p)
	}

	// Register the generator
	engine.RegisterActor(g)

	fmt.Printf("Cores:%v\tservice_rate:%v\tinterarrival_rate:%v\n", cores, mu, lambda)
	engine.Run(duration)
}
