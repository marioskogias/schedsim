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
		var system string
		var duration int
	*/
	var mu = flag.Float64("mu", 0.02, "mu service rate") // default 50usec
	var lambda = flag.Float64("lambda", 0.005, "lambda poisson interarrival")
	var system = flag.String("system", "rtc", "ps or rtc")
	var duration = flag.Float64("duration", 100000000, "experiment duration")
	var quantum = flag.Float64("quantum", 0.5, "processor quantum")

	flag.Parse()

	engine.InitSim()

	//Add a deterministic generator
	//generator := blocks.NewDDGenerator(2, 1)

	//Add an MD generator
	//generator := blocks.NewMDGenerator(0.5, 1)

	//Add an MM generator
	generator := blocks.NewMMGenerator(*lambda, *mu) // 50usec sevice time, lambda 0.005

	//Add a fifo queue
	q := blocks.NewQueue()

	generator.SetOutQueue(q)

	//Init the statistics
	stats := blocks.NewBookKeeper()
	engine.InitStats(stats)

	// FIXME: handle processor type properly
	if *system == "rtc" {
		//Add a run to completion processor
		processor := &blocks.RTCProcessor{}
		processor.SetInQueue(q)
		processor.SetReqDrain(stats)
		engine.RegisterActor(processor)
	} else if *system == "ps" {
		//Add a shared processor
		processor := blocks.NewSharedProcessor(*quantum)
		processor.SetInQueue(q)
		processor.SetReqDrain(stats)
		engine.RegisterActor(processor)
	}

	//Register actors
	engine.RegisterActor(generator)

	fmt.Printf("rho = %v\n", *lambda / *mu)
	//Run till 100000 time units
	engine.Run(*duration)
}
