package internal

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"

	assets "github.com/madeinheaven91/masterbeat/assets"
)

var (
	_, Debug     = os.LookupEnv("MASTERBEAT_DEBUG")
	MsgDump, _   = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	DebugDump, _ = os.OpenFile("debug.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)

	BeatSymbols = make(map[string]string)
)

// a wrapper to work around the mp3.Decode accepting io.ReadCloser
// FIXME: maybe im stupid
type ReaderWrapper struct {
	*bytes.Reader
}

func (rd ReaderWrapper) Close() error {
	return nil
}

func LoadSound(file string) (*beep.Buffer, beep.Format, error) {
	f, err := assets.Assets.ReadFile(file)
	if err != nil {
		return nil, beep.Format{}, err
	}

	reader := ReaderWrapper{bytes.NewReader(f)}

	streamer, format, err := mp3.Decode(reader)
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
	} else {
		BeatSymbols["past"] = ""
		BeatSymbols["current"] = ""
		BeatSymbols["future"] = ""
	}
}
