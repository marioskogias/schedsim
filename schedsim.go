package main

import (
	"flag"
	"fmt"

	"github.com/marioskogias/schedsim/blocks"
	"github.com/marioskogias/schedsim/engine"
)

func main() {

	/*
		var mu, lambda float64
		var processor string
		var duration int
	*/
	var mu = flag.Float64("mu", 0.02, "mu service rate") // default 50usec
	var lambda = flag.Float64("lambda", 0.005, "lambda poisson interarrival")
	var processorType = flag.String("processorType", "rtc", "ts or rtc")
	var duration = flag.Float64("duration", 100000000, "experiment duration")
	var quantum = flag.Float64("quantum", 0.5, "processor quantum")
	var service = flag.String("service", "m", "m or d or lg")

	flag.Parse()

	engine.InitSim()

	//Init the statistics
	stats := blocks.NewBookKeeper()
	engine.InitStats(stats)

	// Create a generator
	var generator engine.ActorInterface
	if *service == "d" {
		//Add an MD generator
		generator = blocks.NewMDGenerator(*lambda, 1 / *mu)
	} else if *service == "m" {
		//Add an MM generator
		generator = blocks.NewMMGenerator(*lambda, *mu) // 50usec sevice time, lambda 0.005
	} else if *service == "lg" {
		// FIXME: make this parametrizable
		// for mean ~ 50 mu = 1 sigma = 2.41
		generator = blocks.NewMLNGenerator(*lambda, 1, 2.41)
	}

	// Create a processor
	var processor blocks.Processor
	if *processorType == "rtc" {
		//Add a run to completion processor
		processor = &blocks.RTCProcessor{}
	} else if *processorType == "ts" {
		//Add a shared processor
		processor = blocks.NewTSProcessor(*quantum)
	} else if *processorType == "ps" {
		processor = blocks.NewPSProcessor()
	}
	processor.SetReqDrain(stats)

	//Add a fifo queue
	q := blocks.NewQueue()

	// Create the topology
	generator.SetOutQueue(q)
	processor.SetInQueue(q)

	// Register Actors
	engine.RegisterActor(generator)
	engine.RegisterActor(processor)

	fmt.Printf("rho = %v\n", *lambda / *mu)
	engine.Run(*duration)
}
