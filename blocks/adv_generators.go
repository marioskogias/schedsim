package blocks

import (
	"bufio"
	"container/list"
	//"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// PBGenerator implements a playback generator for given service times.
// The interarrival distribution is exponential
type PBGenerator struct {
	genericGenerator
	sTimes   [][]int
	cpuCount int
	WaitTime randDist
}

// NewPBGenerator returns a PBGenerator
// Parameters: lambda for the exponential interarrival and the filenames
// with the service times
func NewPBGenerator(lambda float64, paths []string) *PBGenerator {
	g := PBGenerator{}

	for _, p := range paths {
		/* Read service times */
		inFile, _ := os.Open(p)
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		newTimes := make([]int, 0)
		for scanner.Scan() {
			n, _ := strconv.Atoi(scanner.Text())
			newTimes = append(newTimes, n)
		}
		g.sTimes = append(g.sTimes, newTimes)
	}
	g.cpuCount = len(paths)
	g.WaitTime = newExponDistr(lambda)
	return &g
}

// Run is the main loop of the generator
func (g *PBGenerator) Run() {
	for {
		i := rand.Intn(g.cpuCount)
		j := rand.Intn(len(g.sTimes[i]))
		serviceTime := g.sTimes[i][j]
		req := g.Creator.NewRequest(float64(serviceTime))
		g.WriteOutQueueI(req, i)
		g.Wait(g.WaitTime.getRand())
	}
}

// Special generator for Cown scheduling
type Cown struct {
	isSchedulled bool
	queue        *list.List
}

func (c *Cown) GetDelay() float64 {
	panic("Get Delay in cown")
}

func (c *Cown) GetServiceTime() float64 {
	panic("Get ServiceTime in cown")
}

func (c *Cown) SubServiceTime(_ float64) {
	panic("Sub ServiceTime in cown")
}

type cSelRand struct {
	cownCount int
}

func (c cSelRand) getRand() float64 {
	return float64(rand.Intn(c.cownCount))
}

type cSelZipf struct {
	z *rand.Zipf
}

func NewCSelZipf(count int) *cSelZipf {
	c := &cSelZipf{}
	// Initiate Zipf
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c.z = rand.NewZipf(r, 2, 5, uint64(count-1))

	return c
}

func (c *cSelZipf) getRand() float64 {
	return float64(c.z.Uint64())
}

type VeronaGenerator struct {
	genericGenerator
	cownSelector randDist
	cowns        []*Cown
}

func (g *VeronaGenerator) Run() {
	for {
		req := g.Creator.NewRequest(g.ServiceTime.getRand())
		// Assume random cown selection for starters
		cownIdx := int(g.cownSelector.getRand())
		//fmt.Println(cownIdx)
		g.cowns[cownIdx].queue.PushBack(req)
		if !g.cowns[cownIdx].isSchedulled {
			qIdx := rand.Intn(g.GetOutQueueCount())
			g.WriteOutQueueI(g.cowns[cownIdx], qIdx)
			g.cowns[cownIdx].isSchedulled = true
		}
		g.Wait(g.WaitTime.getRand())
	}
}

func newVeronaGenerator(waitLambda float64, cownCount, sel int) *VeronaGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &VeronaGenerator{}
	g.cowns = make([]*Cown, cownCount)
	for i := 0; i < cownCount; i++ {
		g.cowns[i] = &Cown{}
		g.cowns[i].queue = list.New()
	}
	g.WaitTime = newExponDistr(waitLambda)

	if sel == 0 {
		selector := &cSelRand{cownCount}
		g.cownSelector = selector
	} else {
		selector := NewCSelZipf(cownCount)
		g.cownSelector = selector
	}

	return g
}

func NewVeronaMDGenerator(waitLambda, serviceTime float64, cownCount, sel int) *VeronaGenerator {
	g := newVeronaGenerator(waitLambda, cownCount, sel)
	g.ServiceTime = newDeterministicDistr(serviceTime)

	return g
}

func NewVeronaMMGenerator(waitLambda, serviceMu float64, cownCount, sel int) *VeronaGenerator {
	g := newVeronaGenerator(waitLambda, cownCount, sel)
	g.ServiceTime = newExponDistr(serviceMu)

	return g
}

func NewVeronaMBGenerator(waitLambda, peak1, peak2, ratio float64, cownCount, sel int) *VeronaGenerator {
	g := newVeronaGenerator(waitLambda, cownCount, sel)
	g.ServiceTime = newBiDistr(peak1, peak2, ratio)

	return g
}
