package main

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.scanning && !m.scanComplete {
				if m.form.focusIndex == 2 {
					return m.startScan()
				}
				m.form.focusIndex++
				if m.form.focusIndex > 2 {
					m.form.focusIndex = 0
				}
			}
		case tea.KeyTab:
			if !m.scanning && !m.scanComplete {
				m.form.focusIndex++
				if m.form.focusIndex > 2 {
					m.form.focusIndex = 0
				}
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
		switch m.form.focusIndex {
		case 0:
			m.form.startPort, cmd = m.form.startPort.Update(msg)
		case 1:
			m.form.endPort, cmd = m.form.endPort.Update(msg)
		case 2:
			m.form.showFree, cmd = m.form.showFree.Update(msg)
		}
	}

	return m, cmd
}

func (m model) startScan() (tea.Model, tea.Cmd) {
	var err error
	m.startPort, err = strconv.Atoi(strings.TrimSpace(m.form.startPort.Value()))
	if err != nil {
		m.err = err
		return m, nil
	}

	m.endPort, err = strconv.Atoi(strings.TrimSpace(m.form.endPort.Value()))
	if err != nil {
		m.err = err
		return m, nil
	}

	m.showFree = strings.ToLower(strings.TrimSpace(m.form.showFree.Value())) == "y"

	m.scanning = true
	m.err = nil
	return m, scanPort(m.startPort)
}
