package blocks

import (
	"fmt"
	"math"

	"github.com/marioskogias/schedsim/engine"
)

const (
	bUCKET_COUNT = 100000
	gRANULARITY  = 10
)

type Request struct {
	InitTime       float64
	ServiceTime    float64
	serviceTimeImm float64 // This is an immutable version of the service time
	Processed      bool
	DeadLine       float64
	PropDelay      float64
	QoS            int
}

func NewRequest(serviceTime float64) Request {
	return Request{InitTime: engine.GetTime(), ServiceTime: serviceTime, serviceTimeImm: serviceTime}
}

func (r *Request) GetInitialServiceTime() float64 {
	return r.serviceTimeImm
}

func (r *Request) getDelay() float64 {
	return engine.GetTime() - r.InitTime + r.PropDelay
}

func (r Request) GetCmpVal() float64 {
	return r.InitTime
	//d := r.DeadLine - engine.GetTime()
	//return d
}

func (r Request) GetServiceTime() float64 {
	return r.ServiceTime
}

type histogram struct {
	granularity float64
	buckets     []int
	count       int64
	minBucket   int
	maxBucket   int
	sum         float64
	sum_square  float64
}

func newHistogram() *histogram {
	return &histogram{
		granularity: gRANULARITY,
		buckets:     make([]int, bUCKET_COUNT),
		minBucket:   bUCKET_COUNT - 1,
		maxBucket:   0,
	}
}

func (hdr *histogram) addSample(s float64) {
	index := int(s / hdr.granularity)
	if index >= bUCKET_COUNT {
		index = bUCKET_COUNT - 1
	}
	if index < 0 || index >= bUCKET_COUNT {
		panic(fmt.Sprintf("Wrong index: %v\n", index))
	}
	hdr.buckets[index]++
	if index > hdr.maxBucket {
		hdr.maxBucket = index
	}
	if index < hdr.minBucket {
		hdr.minBucket = index
	}
	hdr.count++
	hdr.sum += s
	hdr.sum_square += s * s
}

func (hdr *histogram) avg() float64 {
	return hdr.sum / float64(hdr.count)
}

func (hdr *histogram) stddev() float64 {
	square_avg := hdr.sum_square / float64(hdr.count)
	mean := hdr.avg()

	return math.Sqrt(square_avg - mean*mean)
}

//FIXME: I assume that in every bucket there will be max one percentile
func (hdr *histogram) getPercentiles() map[float64]float64 {
	accum := make([]int, bUCKET_COUNT)
	res := map[float64]float64{}
	percentiles := []float64{0.5, 0.9, 0.95, 0.99}
	percentile_i := 0

	accum[hdr.minBucket] = hdr.buckets[hdr.minBucket]

	// what if percentiles in the first bucket
	for float64(accum[hdr.minBucket]) > percentiles[percentile_i]*float64(hdr.count) {
		// linear interpolation

		res[percentiles[percentile_i]] = hdr.granularity / float64(hdr.buckets[hdr.minBucket]) * (percentiles[percentile_i] * float64(hdr.count))

		percentile_i++
	}
	if percentile_i >= len(percentiles) {
		return res
	}

	for i := hdr.minBucket + 1; i <= hdr.maxBucket; i++ {
		accum[i] = accum[i-1] + hdr.buckets[i]
		for float64(accum[i]) > percentiles[percentile_i]*float64(hdr.count) {
			// linear interpolation
			down := hdr.granularity * float64(i-1)

			res[percentiles[percentile_i]] = down + hdr.granularity/float64(hdr.buckets[i])*(percentiles[percentile_i]*float64(hdr.count)-float64(accum[i-1]))
			percentile_i++
			if percentile_i >= len(percentiles) {
				return res
			}
		}
	}
	return res
}

func (hdr *histogram) printPercentiles() {
	percentiles := hdr.getPercentiles()
	vals := []float64{0.5, 0.9, 0.95, 0.99}
	for _, v := range vals {
		fmt.Printf("%vth: %v\t", int(v*100.0), percentiles[v])
	}
	fmt.Println()

	fmt.Printf("Req/time_unit:%v\n", float64(hdr.count)/engine.GetTime())
}

type BookKeeper struct {
	hdr  *histogram
	name string
}

func NewBookKeeper() *BookKeeper {
	return &BookKeeper{
		hdr: newHistogram(),
	}
}

func (b *BookKeeper) SetName(name string) {
	b.name = name
}

func (b *BookKeeper) TerminateReq(r Request) {
	// FIXME: there is something wrong here
	// panics sometimes with: index out of range
	d := r.getDelay()
	//fmt.Printf("%v\n", d)
	if d < 0 {
		panic("Request with negative service time")
	}
	b.hdr.addSample(d)
}

func (b *BookKeeper) PrintStats() {
	fmt.Printf("%v\nCount: %v AVG: %v STDDev: %v \n", b.name, b.hdr.count, b.hdr.avg(), b.hdr.stddev())
	b.hdr.printPercentiles()
}
