package blocks

import (
	"math/rand"
	"time"

	"github.com/marioskogias/schedsim/engine"
)

type RandDist interface {
	GetRand() float64
}

type genericGenerator struct {
	engine.Actor
	ServiceTime RandDist
	WaitTime    RandDist
}

func (g *genericGenerator) GetGenericActor() *engine.Actor {
	return &g.Actor
}

type RandGenerator struct {
	genericGenerator
}

func (g *RandGenerator) Run() {
	for {
		req := NewRequest(g.ServiceTime.GetRand())
		g.WriteOutQueueI(req, rand.Intn(g.OutQueueCount()))
		g.Wait(g.WaitTime.GetRand())
	}
}

type RRGenerator struct {
	genericGenerator
}

func (g *RRGenerator) Run() {
	for count := 0; ; count++ {
		req := NewRequest(g.ServiceTime.GetRand())
		g.WriteOutQueueI(req, count%g.OutQueueCount())
		g.Wait(g.WaitTime.GetRand())
	}
}

// DDGenerator is a fixed waiting time generator that produces fixed service time requests
type DDGenerator struct {
	RRGenerator
}

func NewDDGenerator(waitTime, serviceTime float64) *DDGenerator {
	g := &DDGenerator{}
	g.ServiceTime = NewDeterministicDistr(serviceTime)
	g.WaitTime = NewDeterministicDistr(waitTime)
	return g
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
type MDGenerator struct {
	RRGenerator
}

func NewMDGenerator(waitLambda float64, serviceTime float64) *MDGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MDGenerator{}
	g.ServiceTime = NewDeterministicDistr(serviceTime)
	g.WaitTime = NewExponDistr(waitLambda)
	return g
}

type MDRandGenerator struct {
	RandGenerator
}

func NewMDRandGenerator(waitLambda float64, serviceTime float64) *MDRandGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MDRandGenerator{}
	g.ServiceTime = NewExponDistr(waitLambda)
	g.WaitTime = NewDeterministicDistr(serviceTime)
	return g
}

// MMGenerator is a exponential waiting time generator that produces exponential service time requests
type MMGenerator struct {
	RRGenerator
}

func NewMMGenerator(waitLambda float64, serviceMu float64) *MMGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MMGenerator{}
	g.ServiceTime = NewExponDistr(serviceMu)
	g.WaitTime = NewExponDistr(waitLambda)
	return g
}

// MMGenerator is a exponential waiting time generator that produces exponential service time requests
type MMRandGenerator struct {
	RandGenerator
}

func NewMMRandGenerator(waitLambda float64, serviceMu float64) *MMRandGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MMRandGenerator{}
	g.ServiceTime = NewExponDistr(serviceMu)
	g.WaitTime = NewExponDistr(waitLambda)
	return g
}

//MLNGenerator is exponential waiting time lognormal service time generator
type MLNGenerator struct {
	RRGenerator
}

func NewMLNGenerator(waitLambda, mu, sigma float64) *MLNGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MLNGenerator{}
	g.ServiceTime = NewLGDistr(mu, sigma)
	g.WaitTime = NewExponDistr(waitLambda)
	return g
}

// DBGenerator (deterministic bimodal) is a poisson interarrival generator with
// requests with bimodal service times (2 values)
type DBGenerator struct {
	RRGenerator
}

func NewDBGenerator(waitLambda, peak1, peak2, ratio float64) *DBGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &DBGenerator{}
	g.ServiceTime = NewBiDistr(peak1, peak2, ratio)
	g.WaitTime = NewExponDistr(waitLambda)
	return g
}
