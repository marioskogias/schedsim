package blocks

import (
	//"fmt"
	"math/rand"
	"time"

	"github.com/marioskogias/schedsim2/engine"
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
