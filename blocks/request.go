package blocks

import (
	"fmt"

	"github.com/marioskogias/schedsim2/engine"
)

type Request struct {
	InitTime    float64
	ServiceTime float64
}

func (r *Request) GetServiceTime() float64 {
	return r.ServiceTime
}

type requestLog struct {
	sum   float64
	count int64
}

func (r *requestLog) addRequest(req Request) {
	r.sum += (engine.GetTime() - req.InitTime)
	r.count += 1
}

func (r *requestLog) avg() float64 {
	return r.sum / float64(r.count)
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
