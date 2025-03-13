package cmd

import (
	"fmt"

	"github.com/gopxl/beep"
)

type MultiTrackSeeker struct {
	tracks   []beep.StreamSeeker
	format   beep.Format
	position int
	length   int
}

func (mts *MultiTrackSeeker) Stream(samples [][2]float64) (n int, ok bool) {
	if mts.position >= mts.length {
		return 0, false
	}

	buffer := make([][2]float64, len(samples))
	for i := range samples {
		samples[i] = [2]float64{0, 0}
	}

	for _, track := range mts.tracks {
		if mts.position < track.Len() {
			nTrack, _ := track.Stream(buffer)
			for i := 0; i < nTrack && i < len(samples); i++ {
				samples[i][0] += buffer[i][0]
				samples[i][1] += buffer[i][1]
			}
		}
	}
	mts.position += len(samples)
	return len(samples), true
}

func (mts *MultiTrackSeeker) Seek(p int) error {
	if p < 0 || p > mts.length {
		return fmt.Errorf("seek position out of range")
	}
	for _, track := range mts.tracks {
		if p < track.Len() {
			if err := track.Seek(p); err != nil {
				return err
			}
		}
	}
	mts.position = p
	return nil
}

func (mts *MultiTrackSeeker) Len() int {
	return mts.length
}

func (mts *MultiTrackSeeker) Position() int {
	return mts.position
}

func NewMultiTrackSeeker(tracks []beep.StreamSeeker, format beep.Format) *MultiTrackSeeker {
	length := 0
	for _, track := range tracks {
		if track.Len() > length {
			length = track.Len()
		}
	}
	return &MultiTrackSeeker{
		tracks:   tracks,
		format:   format,
		position: 0,
		length:   length,
	}
}
