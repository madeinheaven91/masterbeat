package main

import (
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
)

func LoadSound(file string) (*beep.Buffer, beep.Format, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, beep.Format{}, err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, beep.Format{}, err
	}
	defer streamer.Close()

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	return buffer, format, nil
}

func CalculateInterval(bpm float64, sig TimeSignature) time.Duration {
	return 4 * time.Minute / (time.Duration(sig.Bottom) * time.Duration(bpm))
}
