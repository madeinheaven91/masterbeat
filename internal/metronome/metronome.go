package metronome

import (
	"fmt"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"

	misc "github.com/madeinheaven91/masterbeat/internal/misc"
)

type Metronome struct {
	BPM           float64
	Active        bool
	CurrentBeat   int
	Interval      time.Duration
	Sounds        *SoundBank
	ToggleChan    chan struct{}
	TimeSignature *misc.TimeSignature
}

func NewMetronome(bpm float64, timeSignature *misc.TimeSignature, sounds *SoundBank) (*Metronome, error) {
	err := speaker.Init(sounds.format.SampleRate, sounds.format.SampleRate.N(time.Second/50))
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize speaker: %w", err)
	}

	// the higher a bottom number is the faster it goes
	return &Metronome{
		BPM:           bpm,
		Active:        false,
		Interval:      misc.CalculateInterval(bpm, timeSignature),
		Sounds:        sounds,
		ToggleChan:    make(chan struct{}),
		TimeSignature: timeSignature,
	}, nil
}

func (m *Metronome) playSound(sound *beep.Buffer) {
	streamer := sound.Streamer(0, sound.Len())
	speaker.Play(beep.Seq(streamer))
}

func (m *Metronome) Terminate() {
	close(m.ToggleChan)
	speaker.Close()
}

func (m *Metronome) Start() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if misc.Debug {
		spew.Fdump(misc.DebugDump, m.BPM)
	}

	nextBeat := time.Now()
	m.CurrentBeat = 0

	for {
		select {
		case <-m.ToggleChan:
			m.CurrentBeat = 0
			return
		default:
			now := time.Now()
			if now.Before(nextBeat) {
				// Sleep precisely until the next beat
				time.Sleep(nextBeat.Sub(now) - 1*time.Millisecond)
				continue
			}

			if m.CurrentBeat%m.TimeSignature.Top == 0 {
				m.playSound(m.Sounds.accentSound)
				m.CurrentBeat = 0
			} else {
				m.playSound(m.Sounds.regularSound)
			}
			m.CurrentBeat++

			// Calculate next beat time (compensate for any drift)
			nextBeat = nextBeat.Add(m.Interval)

			// Debug output
			// drift := time.Since(nextBeat) - m.interval
			// fmt.Printf("Beat %d (drift: %v)\n", m.beatCount, drift)
		}
	}
}

func (m *Metronome) IncreaseBPM(inc float64) {
	// m.Active = false
	m.BPM += inc
	if m.BPM > 300 {
		m.BPM = 300
	}
	m.Interval = misc.CalculateInterval(m.BPM, m.TimeSignature)
}

func (m *Metronome) DecreaseBPM(dec float64) {
	// m.Active = false
	m.BPM -= dec
	if m.BPM < 30 {
		m.BPM = 30
	}
	m.Interval = misc.CalculateInterval(m.BPM, m.TimeSignature)
}

func (m *Metronome) playSoundCmd(accent bool) tea.Cmd {
	return func() tea.Msg {
		var sound *beep.Buffer
		if accent {
			sound = m.Sounds.accentSound
		} else {
			sound = m.Sounds.regularSound
		}
		streamer := sound.Streamer(0, sound.Len())
		speaker.Play(beep.Seq(streamer))
		return nil
	}
}
