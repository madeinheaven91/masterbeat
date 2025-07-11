package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	bpm := flag.Float64("bpm", 120, "bpm")
	timeSignature := flag.String("time-signature", "4/4", "time signature")
	flag.Parse()

	sounds, err := NewSoundBank("regular.mp3", "accent.mp3")
	if err != nil {
		log.Fatal(err)
	}

	signature, err := ParseTimeSignature(*timeSignature)
	if err != nil {
		log.Fatal(err)
	}

	metronome, err := NewMetronome(*bpm, signature, sounds)
	if err != nil {
		log.Fatal(err)
	}
	defer metronome.Stop()

	fmt.Printf("Starting %d/%d metronome with accent on 4th beat at %.2f bpm\n", signature.Top, signature.Bottom, *bpm)
	fmt.Println("Press Enter to stop...")
	go metronome.Start()

	var input string
	fmt.Scanln(&input)
}
