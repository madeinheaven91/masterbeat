package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"

	"github.com/madeinheaven91/masterbeat/internal/metronome"
	misc "github.com/madeinheaven91/masterbeat/internal/misc"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "00", Dark: "11"}).Bold(true).Underline(true)
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	windowStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Align(lipgloss.Top, lipgloss.Left).
			BorderStyle(lipgloss.ASCIIBorder())
)

type model struct {
	metronome *metronome.Metronome
	status    string
}

type tickMsg time.Time

func InitModel(metronome *metronome.Metronome) *model {
	return &model{
		metronome: metronome,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if misc.Debug {
		spew.Fdump(misc.MsgDump, msg)
	}
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case " ":
			if m.metronome.Active {
				m.metronome.ToggleChan <- struct{}{}
			} else {
				go m.metronome.Start()
			}
			m.metronome.Active = !m.metronome.Active
		case "up":
			// if m.metronome.Active {
			// 	m.metronome.ToggleChan <- struct{}{}
			// }
			m.metronome.IncreaseBPM(1)
		case "shift+up":
			// if m.metronome.Active {
			// 	m.metronome.ToggleChan <- struct{}{}
			// }
			m.metronome.IncreaseBPM(10)
		case "down":
			// if m.metronome.Active {
			// 	m.metronome.ToggleChan <- struct{}{}
			// }
			m.metronome.DecreaseBPM(1)
		case "shift+down":
			// if m.metronome.Active {
			// 	m.metronome.ToggleChan <- struct{}{}
			// }
			m.metronome.DecreaseBPM(10)
		default:
			m.status = msg.String()
		}

	case tickMsg:
		if !m.metronome.Active {
			return m, nil
		}
	}

	return m, tickCmd(m.metronome)
}

func (m model) View() string {
	var metronomeStatus string
	if m.metronome.Active {
		metronomeStatus = "PLAYING"
	} else {
		metronomeStatus = "STOPPED"
	}
	return windowStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		headerStyle.Render("Masterbeat"),
		metronomeStatus,
		BeatsView(m.metronome),
		fmt.Sprintf("%.2f BPM, %s", m.metronome.BPM, m.metronome.TimeSignature),
	))
}

func tickCmd(m *metronome.Metronome) tea.Cmd {
	return tea.Tick(misc.CalculateInterval(m.BPM, m.TimeSignature), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func BeatsView(m *metronome.Metronome) string {
	res := ""
	if m.Active {
		for range m.CurrentBeat - 1 {
			res += dimStyle.Render(misc.BeatSymbols["past"]) + " " 
		}
		res += misc.BeatSymbols["current"] + " "
		for range m.TimeSignature.Top - m.CurrentBeat {
			res += dimStyle.Render(misc.BeatSymbols["future"]) + " "
		}
	} else {
		for range m.TimeSignature.Top {
			res += dimStyle.Render(misc.BeatSymbols["future"]) + " "
		}
	}
	return res
}
