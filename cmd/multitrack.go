package cmd

import (
	"fmt"
	"time"

	"github.com/gopxl/beep"
)

type Track struct {
	Streamer    beep.StreamSeeker
	TrackNumber int
	TrackName   string
	Offset      float64
}

type MultiTrackSeeker struct {
	Tracks   []Track
	format   beep.Format
	position int
	length   int
}

func (mts *MultiTrackSeeker) AddTrackWithOffset(track beep.StreamSeeker, fileName string, offset float64) int {
	nextTrackNumber := 1
	for _, t := range mts.Tracks {
		if t.TrackNumber >= nextTrackNumber {
			nextTrackNumber = t.TrackNumber + 1
		}
	}
	// Create silence streamer for the offset duration.
	silenceSamples := mts.format.SampleRate.N(time.Duration(offset * float64(time.Second)))
	silenceStreamer := beep.Silence(silenceSamples)
	// Combine silence with the actual track using a CompositeSeeker.
	composite := &CompositeSeeker{
		silenceLen: silenceSamples,
		track:      track,
		pos:        0,
	}

	newTrack := Track{
		Streamer:    composite,
		TrackNumber: nextTrackNumber,
		TrackName:   fileName,
		Offset:      offset,
	}
	mts.Tracks = append(mts.Tracks, newTrack)
	// Update overall length: include silence plus track length.
	trackLength := silenceSamples + track.Len()
	if trackLength > mts.length {
		mts.length = trackLength
	}
	return newTrack.TrackNumber
}

func (mts *MultiTrackSeeker) AddTrack(track beep.StreamSeeker, fileName string) int {
	return mts.AddTrackWithOffset(track, fileName, 0)
}

func (mts *MultiTrackSeeker) RemoveTrack(index int) error {
	if index < 0 || index >= len(mts.Tracks) {
		return fmt.Errorf("track index %d out of range", index)
	}
	mts.Tracks = append(mts.Tracks[:index], mts.Tracks[index+1:]...)
	// Recalculate the overall length based on the remaining tracks.
	mts.length = 0
	for _, t := range mts.Tracks {
		if t.Streamer.Len() > mts.length {
			mts.length = t.Streamer.Len()
		}
	}
	return nil
}

func (mts *MultiTrackSeeker) Stream(samples [][2]float64) (n int, ok bool) {
	if mts.position >= mts.length {
		return 0, false
	}

	// Zero out the output buffer.
	for i := range samples {
		samples[i] = [2]float64{0, 0}
	}

	buffer := make([][2]float64, len(samples))
	for _, t := range mts.Tracks {
		nTrack, _ := t.Streamer.Stream(buffer)
		for i := 0; i < nTrack && i < len(samples); i++ {
			samples[i][0] += buffer[i][0]
			samples[i][1] += buffer[i][1]
		}
	}
	mts.position += len(samples)
	return len(samples), true
}

func (mts *MultiTrackSeeker) Seek(p int) error {
	if p < 0 || p > mts.length {
		return fmt.Errorf("seek position out of range")
	}
	for _, t := range mts.Tracks {
		offsetSamples := mts.format.SampleRate.N(time.Duration(t.Offset * float64(time.Second)))
		effectivePos := p - offsetSamples
		if effectivePos < 0 {
			effectivePos = 0
		}
		if effectivePos < t.Streamer.Len() {
			if err := t.Streamer.Seek(effectivePos); err != nil {
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

func (mts *MultiTrackSeeker) Err() error {
	return nil
}

type CompositeSeeker struct {
	silenceLen int
	track      beep.StreamSeeker
	pos        int
}

func (cs *CompositeSeeker) Stream(samples [][2]float64) (n int, ok bool) {
	total := len(samples)
	if cs.pos < cs.silenceLen {
		silenceRemaining := cs.silenceLen - cs.pos
		nSilence := total
		if silenceRemaining < total {
			nSilence = silenceRemaining
		}
		for i := 0; i < nSilence; i++ {
			samples[i][0] = 0
			samples[i][1] = 0
		}
		cs.pos += nSilence
		n += nSilence
		if nSilence < total {
			nTrack, okTrack := cs.track.Stream(samples[nSilence:])
			n += nTrack
			cs.pos += nTrack
			return n, okTrack
		}
		return n, true
	} else {
		nTrack, okTrack := cs.track.Stream(samples)
		cs.pos += nTrack
		return nTrack, okTrack
	}
}

func (cs *CompositeSeeker) Seek(p int) error {
	if p < 0 || p > cs.Len() {
		return fmt.Errorf("seek position out of range")
	}
	cs.pos = p
	if p < cs.silenceLen {
		return cs.track.Seek(0)
	}
	return cs.track.Seek(p - cs.silenceLen)
}

func (cs *CompositeSeeker) Len() int {
	return cs.silenceLen + cs.track.Len()
}

func NewMultiTrackSeeker(streams []beep.StreamSeeker, format beep.Format) *MultiTrackSeeker {
	length := 0
	tracks := make([]Track, 0, len(streams))
	for i, s := range streams {
		if s.Len() > length {
			length = s.Len()
		}
		tracks = append(tracks, Track{
			Streamer:    s,
			TrackNumber: i + 1,
			TrackName:   fmt.Sprintf("Track %d", i+1),
		})
	}
	return &MultiTrackSeeker{
		Tracks:   tracks,
		format:   format,
		position: 0,
		length:   length,
	}
}
