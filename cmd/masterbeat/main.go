package main

import (
	"flag"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/madeinheaven91/masterbeat/internal/metronome"
	misc "github.com/madeinheaven91/masterbeat/internal/misc"
	"github.com/madeinheaven91/masterbeat/internal/ui"
)

func main() {
	bpm := flag.Float64("bpm", 120, "bpm")
	timeSignature := flag.String("time-signature", "4/4", "time signature")
	noNerd := flag.Bool("no-nerd", false, "use this if you don't use nerd fonts")
	flag.Parse()

	signature, err := misc.ParseTimeSignature(*timeSignature)
	if err != nil {
		misc.Error("Time signature parsing error: %v\n", err)
		return
	}

	sounds, err := metronome.NewSoundBank("regular.mp3", "accent.mp3")
	if err != nil {
		misc.Error("Soundbank loading error: %v\n", err)
		return
	}

	metronome, err := metronome.NewMetronome(*bpm, signature, sounds)
	if err != nil {
		misc.Error("Metronome initializing error: %v\n", err)
		return
	}
	defer metronome.Terminate()

	misc.InitBeatSymbols(*noNerd)

	model := ui.InitModel(metronome)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		misc.Error("An error occured: %v\n", err)
		os.Exit(1)
	}
}
