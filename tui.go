package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var m = initialModel()
var p = tea.NewProgram(m)

type (
	errMsg error
	logMsg Log
)

type LogType int

const (
	Info LogType = iota
	Warning
	Error
)

type Log struct {
	Msg  string
	Type LogType
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

var (
	InfoLogStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66F359"))
	WarningLogStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F3F359"))
	ErrorLogStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F35959"))
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(80)
	ta.SetHeight(2)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(80, 10)
	vp.SetContent("Type a prompt and press Enter to ask Clyde AI.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			prompt := m.textarea.Value()

			go SendDiscordMessage(prompt)

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+prompt)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case errMsg:
		m.err = msg
		return m, nil

	case DiscordMessage:
		m.messages = append(m.messages, m.senderStyle.Render("AI: ")+msg.Content)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

		return m, nil

	case logMsg:
		switch msg.Type {
		case Info:
			m.messages = append(m.messages, InfoLogStyle.Render("SYSTEM: "+msg.Msg))
		case Warning:
			m.messages = append(m.messages, WarningLogStyle.Render("SYSTEM: "+msg.Msg))
		case Error:
			m.messages = append(m.messages, ErrorLogStyle.Render("SYSTEM: "+msg.Msg))
		}

		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
