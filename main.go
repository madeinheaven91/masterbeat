package main

import (
	"flag"
	"fmt"
)

func main() {
	bpm := flag.Float64("bpm", 120, "bpm")
	timeSignature := flag.String("time-signature", "4/4", "time signature")
	flag.Parse()

	sounds, err := NewSoundBank("regular.mp3", "accent.mp3")
	if err != nil {
		Error("Soundbank loading error: %s\n", err)
		return
	}

	signature, err := ParseTimeSignature(*timeSignature)
	if err != nil {
		Error("Time signature parsing error: %s\n", err)
		return
	}

	metronome, err := NewMetronome(*bpm, signature, sounds)
	if err != nil {
		Error("Metronome initializing error: %s\n", err)
		return
	}
	defer metronome.Stop()

	fmt.Printf("Starting %d/%d metronome at %.2f bpm\n", signature.Top, signature.Bottom, *bpm)
	fmt.Println("Press Enter to stop...")
	go metronome.Start()

	var input string
	fmt.Scanln(&input)
}
