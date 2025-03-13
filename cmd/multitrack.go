package cmd

import (
	"fmt"

	"github.com/gopxl/beep"
)

type Track struct {
	Streamer    beep.StreamSeeker
	TrackNumber int
	TrackName   string
}

type MultiTrackSeeker struct {
	Tracks   []Track
	format   beep.Format
	position int
	length   int
}

func (mts *MultiTrackSeeker) AddTrack(track beep.StreamSeeker, fileName string) int {
	nextTrackNumber := 1
	for _, t := range mts.Tracks {
		if t.TrackNumber >= nextTrackNumber {
			nextTrackNumber = t.TrackNumber + 1
		}
	}
	newTrack := Track{
		Streamer:    track,
		TrackNumber: nextTrackNumber,
		TrackName:   fileName,
	}
	mts.Tracks = append(mts.Tracks, newTrack)
	// Update the overall length if the added track is longer.
	if track.Len() > mts.length {
		mts.length = track.Len()
	}
	return newTrack.TrackNumber
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

	buffer := make([][2]float64, len(samples))
	for i := range samples {
		samples[i] = [2]float64{0, 0}
	}

	for _, t := range mts.Tracks {
		if mts.position < t.Streamer.Len() {
			nTrack, _ := t.Streamer.Stream(buffer)
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
	for _, t := range mts.Tracks {
		if p < t.Streamer.Len() {
			if err := t.Streamer.Seek(p); err != nil {
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
