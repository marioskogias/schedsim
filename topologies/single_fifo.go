package topologies

import (
	"fmt"

	"github.com/marioskogias/schedsim/blocks"
	"github.com/marioskogias/schedsim/engine"
)

func SingleQueue(lambda, mu, duration float64) {

	engine.InitSim()

	//Init the statistics
	stats := blocks.NewBookKeeper()
	engine.InitStats(stats)

	// Add generator
	//g := blocks.NewMMGenerator(lambda, mu) // 50usec sevice time, lambda 0.005
	g := blocks.NewMDGenerator(lambda, 10)

	// Create queues
	q := blocks.NewQueue()

	// Create processors
	processors := make([]blocks.Processor, cores)

	// first the slow cores
	for i := 0; i < cores; i++ {
		processors[i] = &blocks.RTCProcessor{}
	}

	// Connect the queue
	g.AddOutQueue(q)

	for i := 0; i < cores; i++ {
		processors[i].AddInQueue(q)
	}

	// Add the stats and register processors
	for _, p := range processors {
		p.SetReqDrain(stats)
		engine.RegisterActor(p)
	}

	// Register the generator
	engine.RegisterActor(g)

	fmt.Printf("rho=%v\toffered_qps=%v\tservice_rate=%v\n", lambda/mu, lambda, mu)
	engine.Run(duration)
}
