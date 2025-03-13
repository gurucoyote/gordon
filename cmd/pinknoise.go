package cmd

import (
	"math/bits"
	"math/rand"
	"time"

	"github.com/gopxl/beep"
)

// PinkNoise implements beep.Streamer and generates pink noise using the Voss–McCartney algorithm.
type PinkNoise struct {
	rng        *rand.Rand
	rows       []float64
	runningSum float64
	index      uint64
	numRows    int
}

// NewPinkNoise creates a new PinkNoise streamer.
func NewPinkNoise() *PinkNoise {
	numRows := 16
	rows := make([]float64, numRows)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sum float64
	for i := 0; i < numRows; i++ {
		rows[i] = rng.Float64()*2 - 1
		sum += rows[i]
	}

	return &PinkNoise{
		rng:        rng,
		rows:       rows,
		runningSum: sum,
		index:      0,
		numRows:    numRows,
	}
}

// Stream fills the provided samples slice with pink noise samples.
func (p *PinkNoise) Stream(samples [][2]float64) (n int, ok bool) {
	const scale = 1.0
	for i := range samples {
		p.index++
		zeros := bits.TrailingZeros64(p.index)
		if zeros < p.numRows {
			oldVal := p.rows[zeros]
			newVal := p.rng.Float64()*2 - 1
			p.rows[zeros] = newVal
			p.runningSum += newVal - oldVal
		}
		sampleValue := scale * (p.runningSum / float64(p.numRows))
		samples[i][0] = sampleValue
		samples[i][1] = sampleValue
		n++
	}
	return n, true
}

// Err returns nil as this Streamer does not track errors.
func (p *PinkNoise) Err() error {
	return nil
}
