package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/gopxl/beep/v2/mp3"
)

var (
	_, Debug = os.LookupEnv("DEBUG")
	MsgDump, _  = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	DebugDump, _ = os.OpenFile("debug.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)

	BeatSymbols = make(map[string]string)
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

func CalculateInterval(bpm float64, sig *TimeSignature) time.Duration {
	return 4 * time.Minute / (time.Duration(sig.Bottom) * time.Duration(bpm))
}

func Error(format string, a ...any) {
	fmt.Printf("\033[31m[ERROR]\033[0m ")
	fmt.Printf(format, a...)
}

func InitBeatSymbols(noNerd bool) {
	if Debug {
		spew.Fdump(DebugDump, noNerd)
	}
	if noNerd {
		BeatSymbols["past"] = "#"
		BeatSymbols["current"] = "#"
		BeatSymbols["future"] = "-"
	}else {
		BeatSymbols["past"] = ""
		BeatSymbols["current"] = ""
		BeatSymbols["future"] = ""
	}
}
