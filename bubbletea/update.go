package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyTab:
			if !m.scanning && !m.scanComplete {
				// Only start scan if we're on the last input and all fields are filled
				if m.form.focusIndex == 2 && m.allFieldsFilled() {
					return m.startScan()
				}
				// Move to next input field
				m.form.focusIndex = (m.form.focusIndex + 1) % 3
				return m.updateFocus()
			}
		case tea.KeyShiftTab:
			if !m.scanning && !m.scanComplete {
				// Move to previous input field
				m.form.focusIndex = (m.form.focusIndex - 1 + 3) % 3
				return m.updateFocus()
			}
		}

	case portStatus:
		m.results[msg.port] = msg.info
		m.progress++
		if m.progress >= m.endPort-m.startPort+1 {
			m.scanning = false
			m.scanComplete = true
			return m, nil
		}
		return m, scanPort(m.startPort + m.progress)
	}

	if !m.scanning && !m.scanComplete {
		m.form.startPort, cmd = m.form.startPort.Update(msg)
		m.form.endPort, _ = m.form.endPort.Update(msg)
		m.form.showFree, _ = m.form.showFree.Update(msg)
	}

	return m, cmd
}

func (m model) updateFocus() (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 3)
	m.form.startPort.Blur()
	m.form.endPort.Blur()
	m.form.showFree.Blur()
	switch m.form.focusIndex {
	case 0:
		m.form.startPort.Focus()
		cmds[0] = textinput.Blink
	case 1:
		m.form.endPort.Focus()
		cmds[1] = textinput.Blink
	case 2:
		m.form.showFree.Focus()
		cmds[2] = textinput.Blink
	}
	return m, tea.Batch(cmds...)
}

func (m model) allFieldsFilled() bool {
	return m.form.startPort.Value() != "" &&
		m.form.endPort.Value() != "" &&
		m.form.showFree.Value() != ""
}

func (m model) startScan() (tea.Model, tea.Cmd) {
	var err error
	m.startPort, err = strconv.Atoi(strings.TrimSpace(m.form.startPort.Value()))
	if err != nil {
		m.err = fmt.Errorf("invalid start port: %v", err)
		return m, nil
	}

	m.endPort, err = strconv.Atoi(strings.TrimSpace(m.form.endPort.Value()))
	if err != nil {
		m.err = fmt.Errorf("invalid end port: %v", err)
		return m, nil
	}

	if m.startPort > m.endPort {
		m.err = fmt.Errorf("start port must be less than or equal to end port")
		return m, nil
	}

	showFree := strings.ToLower(strings.TrimSpace(m.form.showFree.Value()))
	if showFree != "y" && showFree != "n" {
		m.err = fmt.Errorf("please enter 'y' or 'n' for show free ports")
		return m, nil
	}

	m.showFree = showFree == "y"
	m.scanning = true
	m.err = nil
	return m, scanPort(m.startPort)
}
