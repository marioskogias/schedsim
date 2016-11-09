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
	InitTime    float64
	ServiceTime float64
}

func (r *Request) GetServiceTime() float64 {
	return r.ServiceTime
}

func (r *Request) getDelay() float64 {
	return engine.GetTime() - r.InitTime
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
	if index > bUCKET_COUNT {
		index = bUCKET_COUNT - 1
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
	// i assume only 50th can be in the first bucket -> FIXME
	if float64(accum[0]) > percentiles[percentile_i]*float64(hdr.count) {
		// linear interpolation

		res[percentiles[percentile_i]] = hdr.granularity / float64(hdr.buckets[0]) * (percentiles[percentile_i] * float64(hdr.count))

		percentile_i++
	}

	for i := hdr.minBucket + 1; i <= hdr.maxBucket; i++ {
		accum[i] = accum[i-1] + hdr.buckets[i]

		if float64(accum[i]) > percentiles[percentile_i]*float64(hdr.count) {
			// linear interpolation
			down := hdr.granularity * float64(i-1)

			res[percentiles[percentile_i]] = down + hdr.granularity/float64(hdr.buckets[i])*(percentiles[percentile_i]*float64(hdr.count)-float64(accum[i-1]))
			percentile_i++
			if percentile_i >= len(percentiles) {
				break
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
}

type BookKeeper struct {
	hdr *histogram
}

func NewBookKeeper() *BookKeeper {
	return &BookKeeper{
		hdr: newHistogram(),
	}
}

func (b *BookKeeper) TerminateReq(r Request) {
	// FIXME: there is something wrong here
	// panics sometimes with: index out of range
	d := r.getDelay()
	//fmt.Printf("%v\n", d)
	if d < 0 {
		panic("Request with negative service time")
	}
	b.hdr.addSample(r.getDelay())
}

func (b *BookKeeper) PrintStats() {
	fmt.Printf("Count: %v AVG: %v STDDev: %v \n", b.hdr.count, b.hdr.avg(), b.hdr.stddev())
	b.hdr.printPercentiles()
}
