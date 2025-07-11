package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Metronome struct {
	bpm           float64
	beatCount     int
	interval      time.Duration
	sounds        *SoundBank
	stopChan      chan struct{}
	timeSignature *TimeSignature
}

func NewMetronome(bpm float64, timeSignature *TimeSignature, sounds *SoundBank) (*Metronome, error) {
	err := speaker.Init(sounds.format.SampleRate, sounds.format.SampleRate.N(time.Second/50))
	if err != nil {
		return nil, fmt.Errorf("initializing speaker: %w", err)
	}

	// the higher denominator the faster it goes
	return &Metronome{
		bpm:           bpm,
		interval:      CalculateInterval(bpm, *timeSignature),
		sounds:        sounds,
		stopChan:      make(chan struct{}),
		timeSignature: timeSignature,
	}, nil
}

func (m *Metronome) playSound(sound *beep.Buffer) {
	streamer := sound.Streamer(0, sound.Len())
	speaker.Play(beep.Seq(streamer))
}

func (m *Metronome) Start() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	nextBeat := time.Now()
	m.beatCount = 0

	for {
		select {
		case <-m.stopChan:
			return
		default:
			now := time.Now()
			if now.Before(nextBeat) {
				// Sleep precisely until the next beat
				time.Sleep(nextBeat.Sub(now) - 1*time.Millisecond)
				continue
			}

			// Determine which sound to play
			m.beatCount++
			if m.beatCount%m.timeSignature.Top == 1 {
				m.playSound(m.sounds.accentSound)
			} else {
				m.playSound(m.sounds.regularSound)
			}

			// Calculate next beat time (compensate for any drift)
			nextBeat = nextBeat.Add(m.interval)

			// Debug output
			// drift := time.Since(nextBeat) - m.interval
			// fmt.Printf("Beat %d (drift: %v)\n", m.beatCount, drift)
		}
	}
}

func (m *Metronome) Stop() {
	close(m.stopChan)
	speaker.Close()
}
