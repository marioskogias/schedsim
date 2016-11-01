package blocks

import (
	"fmt"

	"github.com/marioskogias/schedsim2/engine"
)

type Request struct {
	InitTime    int
	ServiceTime int
}

func (r *Request) GetServiceTime() int {
	return r.ServiceTime
}

type requestLog struct {
	sum   int
	count int
}

func (r *requestLog) addRequest(req Request) {
	r.sum += (engine.GetTime() - req.InitTime)
	r.count += 1
}

func (r *requestLog) avg() float64 {
	return float64(r.sum) / float64(r.count)
}

type BookKeeper struct {
	log requestLog
}

func (b *BookKeeper) TerminateReq(r Request) {
	b.log.addRequest(r)
}

func (b *BookKeeper) PrintStats() {
	fmt.Printf("The stats are:\n")
	fmt.Printf("AVG: %v\n", b.log.avg())
}
