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

	//Add a fifo queue
	q := blocks.NewQueue()

	if *service == "d" {
		//Add an MD generator
		generator := blocks.NewMDGenerator(*lambda, 1 / *mu)
		generator.SetOutQueue(q)
		engine.RegisterActor(generator)

	} else if *service == "m" {
		//Add an MM generator
		generator := blocks.NewMMGenerator(*lambda, *mu) // 50usec sevice time, lambda 0.005
		generator.SetOutQueue(q)
		engine.RegisterActor(generator)
	} else if *service == "lg" {
		// FIXME: make this parametrizable
		// for mean ~ 50 mu = 1 sigma = 2.41
		generator := blocks.NewMLNGenerator(*lambda, 1, 2.41)
		generator.SetOutQueue(q)
		engine.RegisterActor(generator)
	}

	//Add a deterministic generator
	//generator := blocks.NewDDGenerator(2, 1)

	//Init the statistics
	stats := blocks.NewBookKeeper()
	engine.InitStats(stats)

	var processor blocks.Processor
	// FIXME: handle processor type properly
	if *processorType == "rtc" {
		//Add a run to completion processor
		processor = &blocks.RTCProcessor{}
	} else if *processorType == "ts" {
		//Add a shared processor
		processor = blocks.NewTSProcessor(*quantum)
	} else if *processorType == "ps" {
		processor = blocks.NewPSProcessor()
	}
	processor.SetInQueue(q)
	processor.SetReqDrain(stats)
	engine.RegisterActor(processor)

	//Register actors

	fmt.Printf("rho = %v\n", *lambda / *mu)
	//Run till 100000 time units
	engine.Run(*duration)
}
