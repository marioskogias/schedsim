package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/epfl-dcsl/schedsim/topologies"
)

func main() {
	var topo = flag.Int("topo", 0, "topology selector")
	var mu = flag.Float64("mu", 0.02, "mu service rate") // default 50usec
	var lambda = flag.String("lambda", "0.005", "lambda poisson interarrival")
	var genType = flag.Int("genType", 0, "type of generator")
	var procType = flag.Int("procType", 0, "type of processor")
	var duration = flag.Float64("duration", 1000000, "experiment duration")
	var bufferSize = flag.Int("buffersize", 1, "size of the bounded buffer")
	var quantum = flag.Float64("quantum", 1.0, "quantum for TS processors")
	var bCount = flag.Int("bCount", 1, "behaviour count per execution")
	var cownCount = flag.Int("cownCount", 1, "total number of cowns")
	var sel = flag.Int("sel", 0, "cown selector: 0 for rand 1 for zipf")
	var fair = flag.Bool("fair", true, "fairness")

	flag.Parse()

	lambdas := strings.Split(*lambda, ":")
	lambdaList := make([]float64, 0)
	if len(lambdas) > 1 {
		start, _ := strconv.ParseFloat(lambdas[0], 64)
		end, _ := strconv.ParseFloat(lambdas[1], 64)
		step, _ := strconv.ParseFloat(lambdas[2], 64)
		for start <= end {
			lambdaList = append(lambdaList, start)
			start += step
		}
	} else {
		val, _ := strconv.ParseFloat(lambdas[0], 64)
		lambdaList = append(lambdaList, val)
	}

	for _, l := range lambdaList {
		if *topo == 0 {
			topologies.SingleQueue(l, *mu, *duration, *genType, *procType)
		} else if *topo == 1 {
			topologies.MultiQueue(l, *mu, *duration, *genType, *procType)
		} else if *topo == 2 {
			topologies.BoundedQueue(l, *mu, *duration, *bufferSize)
		} else if *topo == 3 {
			topologies.Verona(l, *mu, *duration, *genType, *quantum)
		} else if *topo == 4 {
			topologies.VeronaCown(l, *mu, *duration, *genType, *cownCount, *bCount, *sel, *fair)
		} else {
			panic("Unknown topology")
		}
	}
}
