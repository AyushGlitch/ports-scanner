package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputForm struct {
	startPort  textinput.Model
	endPort    textinput.Model
	showFree   textinput.Model
	focusIndex int
}

type model struct {
	form         inputForm
	startPort    int
	endPort      int
	results      map[int]portInfo
	progress     int
	showFree     bool
	scanning     bool
	scanComplete bool
	err          error
}

type portInfo struct {
	isOpen  bool
	service string
}

func initialModel() model {
	m := model{
		form: inputForm{
			startPort:  textinput.New(),
			endPort:    textinput.New(),
			showFree:   textinput.New(),
			focusIndex: 0,
		},
		results: make(map[int]portInfo),
	}

	m.form.startPort.Placeholder = "Start Port"
	m.form.startPort.Focus()
	m.form.startPort.CharLimit = 5
	m.form.startPort.Width = 20

	m.form.endPort.Placeholder = "End Port"
	m.form.endPort.CharLimit = 5
	m.form.endPort.Width = 20

	m.form.showFree.Placeholder = "Show free ports? (y/n)"
	m.form.showFree.CharLimit = 1
	m.form.showFree.Width = 20

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type portStatus struct {
	port int
	info portInfo
}

func scanPort(port int) tea.Cmd {
	return func() tea.Msg {
		address := net.JoinHostPort("localhost", strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
		isOpen := err == nil
		service := ""
		if isOpen {
			service = getServiceName(port)
			conn.Close()
		}
		return portStatus{port: port, info: portInfo{isOpen: isOpen, service: service}}
	}
}

func getServiceName(port int) string {
	switch port {
	case 80:
		return "HTTP"
	case 443:
		return "HTTPS"
	case 22:
		return "SSH"
	case 21:
		return "FTP"
	// Add more well-known ports as needed
	default:
		return fmt.Sprintf("Unknown (%d)", port)
	}
}
