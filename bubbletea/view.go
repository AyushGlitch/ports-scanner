package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

func (m model) View() string {
	var b strings.Builder

	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
		return b.String()
	}

	if !m.scanning && !m.scanComplete {
		b.WriteString("Port Scanner\n\n")
		b.WriteString(m.inputField("Start Port", m.form.startPort, 0))
		b.WriteString(m.inputField("End Port", m.form.endPort, 1))
		b.WriteString(m.inputField("Show free ports? (y/n)", m.form.showFree, 2))
		b.WriteString("\nPress Enter to start scanning\n")
	} else if m.scanning {
		b.WriteString(fmt.Sprintf("Scanning ports %d to %d...\n", m.startPort, m.endPort))
		b.WriteString(fmt.Sprintf("Progress: %d/%d\n\n", m.progress, m.endPort-m.startPort+1))
	} else if m.scanComplete {
		b.WriteString("Scan complete!\n\n")

		portType := "occupied"
		if m.showFree {
			portType = "free"
		}
		b.WriteString(fmt.Sprintf("List of %s ports:\n", portType))

		for port := m.startPort; port <= m.endPort; port++ {
			info, exists := m.results[port]
			if exists && info.isOpen != m.showFree {
				if m.showFree {
					b.WriteString(fmt.Sprintf("- %d\n", port))
				} else {
					b.WriteString(fmt.Sprintf("- %d: %s\n", port, info.service))
				}
			}
		}
	}

	return b.String()
}

func (m model) inputField(label string, input textinput.Model, index int) string {
	if index == m.form.focusIndex {
		return focusedStyle.Render(fmt.Sprintf("%s: %s\n", label, input.View()))
	}
	return blurredStyle.Render(fmt.Sprintf("%s: %s\n", label, input.View()))
}
