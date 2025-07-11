package main

import (
	"fmt"

	"github.com/gopxl/beep/v2"
)

type SoundBank struct {
	regularSound *beep.Buffer
	accentSound  *beep.Buffer
	format       beep.Format
}

func NewSoundBank(regularFile, accentFile string) (*SoundBank, error) {
	regular, format, err := LoadSound(regularFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't load regular sound: %w", err)
	}

	accent, _, err := LoadSound(accentFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't load accent sound: %w", err)
	}

	return &SoundBank{
		regularSound: regular,
		accentSound:  accent,
		format:       format,
	}, nil
}
