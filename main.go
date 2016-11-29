package main

import (
	"flag"
	"fmt"

	"github.com/marioskogias/schedsim/topologies"
)

func main() {
	var topo = flag.Int("topo", 0, "topology selector")
	var mu = flag.Float64("mu", 0.02, "mu service rate") // default 50usec
	var lambda = flag.Float64("lambda", 0.005, "lambda poisson interarrival")
	var duration = flag.Float64("duration", 10000000, "experiment duration")

	flag.Parse()
	fmt.Printf("Selected topology: %v\n", *topo)

	topologies.SingleQueue(*lambda, *mu, *duration)
}
