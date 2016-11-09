package blocks

import (
	//"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/marioskogias/schedsim/engine"
)

type genericGenerator struct {
	engine.Actor
}

func (g *genericGenerator) GetGenericActor() *engine.Actor {
	return &g.Actor
}

// DDGenerator is a fixed waiting time generator that produces fixed service time requests
type DDGenerator struct {
	genericGenerator
	waitTime    float64
	serviceTime float64
}

func NewDDGenerator(waitTime, serviceTime float64) *DDGenerator {
	return &DDGenerator{waitTime: waitTime, serviceTime: serviceTime}
}

func (g *DDGenerator) Run() {
	for {
		//fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		req := NewRequest(g.serviceTime)
		g.WriteOutQueue(req)
		g.Wait(g.waitTime)
	}
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
type MDGenerator struct {
	genericGenerator
	waitLambda  float64
	serviceTime float64
}

func NewMDGenerator(waitLambda float64, serviceTime float64) *MDGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	return &MDGenerator{waitLambda: waitLambda, serviceTime: serviceTime}
}

func (g *MDGenerator) getDelay() float64 {
	d := float64(rand.ExpFloat64() / g.waitLambda)
	//fmt.Printf("%v\n", d)
	return d
}

func (g *MDGenerator) Run() {
	for {
		//fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		req := NewRequest(g.serviceTime)
		g.WriteOutQueue(req)
		g.Wait(g.getDelay())
	}
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
type MMGenerator struct {
	genericGenerator
	waitLambda float64
	serviceMu  float64
}

func NewMMGenerator(waitLambda float64, serviceMu float64) *MMGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	return &MMGenerator{waitLambda: waitLambda, serviceMu: serviceMu}
}

func (g *MMGenerator) getDelay() float64 {
	d := float64(rand.ExpFloat64() / g.waitLambda)
	//fmt.Printf("%v\n", d)
	return d
}

func (g *MMGenerator) getServiceTime() float64 {
	s := float64(rand.ExpFloat64() / g.serviceMu)
	//fmt.Printf("%v\n", s)
	return s
}

func (g *MMGenerator) Run() {
	for {
		//fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		req := NewRequest(g.getServiceTime())
		g.WriteOutQueue(req)
		g.Wait(g.getDelay())
	}
}

//MLNGenerator is exponential waiting time lognormal service time generator
type MLNGenerator struct {
	genericGenerator
	waitLambda float64
	mu         float64
	sigma      float64
}

func NewMLNGenerator(waitLambda, mu, sigma float64) *MLNGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	return &MLNGenerator{waitLambda: waitLambda, mu: mu, sigma: sigma}
}

func (g *MLNGenerator) getDelay() float64 {
	d := float64(rand.ExpFloat64() / g.waitLambda)
	//fmt.Printf("%v\n", d)
	return d
}

func (g *MLNGenerator) getServiceTime() float64 {
	z := rand.NormFloat64()
	s := math.Exp(g.mu + g.sigma*z)
	//fmt.Printf("%v\n", s)
	return s
}

func (g *MLNGenerator) Run() {
	for {
		//fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		req := NewRequest(g.getServiceTime())
		g.WriteOutQueue(req)
		g.Wait(g.getDelay())
	}
}
