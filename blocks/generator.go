package blocks

import (
	//"fmt"
	"math/rand"
	"time"

	"github.com/marioskogias/schedsim/engine"
)

// DDGenerator is a fixed waiting time generator that produces fixed service time requests
type DDGenerator struct {
	engine.Actor
	waitTime    float64
	serviceTime float64
}

func NewDDGenerator(waitTime, serviceTime float64) *DDGenerator {
	return &DDGenerator{waitTime: waitTime, serviceTime: serviceTime}
}

func (g *DDGenerator) Run() {
	for {
		//fmt.Printf("Generator: will add in queue TIME = %v\n", engine.GetTime())
		req := Request{InitTime: engine.GetTime(), ServiceTime: g.serviceTime}
		g.WriteOutQueue(req)
		g.Wait(g.waitTime)
	}
}

func (g *DDGenerator) GetGenericActor() *engine.Actor {
	return &g.Actor
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
type MDGenerator struct {
	engine.Actor
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
		req := Request{InitTime: engine.GetTime(), ServiceTime: g.serviceTime}
		g.WriteOutQueue(req)
		g.Wait(g.getDelay())
	}
}

func (g *MDGenerator) GetGenericActor() *engine.Actor {
	return &g.Actor
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
type MMGenerator struct {
	engine.Actor
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
		req := Request{InitTime: engine.GetTime(), ServiceTime: g.getServiceTime()}
		g.WriteOutQueue(req)
		g.Wait(g.getDelay())
	}
}

func (g *MMGenerator) GetGenericActor() *engine.Actor {
	return &g.Actor
}
