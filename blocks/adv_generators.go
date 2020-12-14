package blocks

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"container/list"
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
	queue *list.List
}

func (c *Cown) GetDelay() float64 {
	panic("Get Delay in cown")
}

func (c *Cown) GetServiceTime() float64 {
	panic("Get ServiceTime in cown")
}

func (c *Cown) SubServiceTime(_ float64)  {
	panic("Sub ServiceTime in cown")
}

type veronaGenerator struct {
	genericGenerator
	cowns []*Cown
}


func (g *veronaGenerator) Run() {
	for {
		req := g.Creator.NewRequest(g.ServiceTime.getRand())
		// Assume random cown selection for starters
		cownIdx := rand.Intn(len(g.cowns))
		g.cowns[cownIdx].queue.PushBack(req)
		if !g.cowns[cownIdx].isSchedulled {
			qIdx := rand.Intn(g.GetOutQueueCount())
			g.WriteOutQueueI(g.cowns[cownIdx], qIdx)
		}
		g.Wait(g.WaitTime.getRand())
	}
}
