package blocks

import (
	"math"
	"math/rand"
)

// Deterministic Distribution
type DeterministicDistr struct {
	d float64
}

func NewDeterministicDistr(d float64) *DeterministicDistr {
	return &DeterministicDistr{d}
}

func (distr *DeterministicDistr) GetRand() float64 {
	return distr.d
}

// Exponential Distribution
type ExponDistr struct {
	lambda float64
}

func NewExponDistr(l float64) *ExponDistr {
	return &ExponDistr{l}
}

func (distr *ExponDistr) GetRand() float64 {
	return float64(rand.ExpFloat64() / distr.lambda)
}

// LogNormal Distribution
type LGDistr struct {
	mu    float64
	sigma float64
}

func NewLGDistr(mu, sigma float64) *LGDistr {
	return &LGDistr{mu, sigma}
}

func (distr *LGDistr) GetRand() float64 {
	z := rand.NormFloat64()
	s := math.Exp(distr.mu + distr.sigma*z)
	return s
}

// Bimodel Distribution
type BiDistr struct {
	v1    float64
	v2    float64
	ratio float64
}

func NewBiDistr(v1, v2, ratio float64) *BiDistr {
	return &BiDistr{v1, v2, ratio}
}

func (distr *BiDistr) GetRand() float64 {
	if rand.Float64() > distr.ratio {
		return distr.v2
	}
	return distr.v1
}
