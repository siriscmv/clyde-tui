package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var tui *tea.Program

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

const (
	spacebar = " "
	help     = "Ctrl+x • Copy last msg | ctrl+s • Toggle mouse | ctrl+p • Multiline prompt | @cb • Paste clipboard"
)

type Log struct {
	Msg  string
	Type LogType
}

type model struct {
	viewport viewport.Model
	textarea textarea.Model
	spinner  spinner.Model
	lastMsg  string
	waiting  bool
	mouse    bool
	messages []string
	err      error
}

type KeyMap struct {
	PageDown     key.Binding
	PageUp       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	Down         key.Binding
	Up           key.Binding
}

var (
	BoldStyle    = lipgloss.NewStyle().Bold(true)
	InfoLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#a6da95"))
	WarningLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#eed49f"))
	ErrorLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#ed8796"))
	UserStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#c6a0f6"))
	FadedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	HelpStyle      = FadedStyle.Copy().Italic(true).Padding(0, 1).Margin(0, 1)
	ContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#c6a0f6")).Padding(1).Margin(1)
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Talk with Clyde here"
	ta.Focus()

	ta.Prompt = UserStyle.Render("❯ ")
	ta.CharLimit = 2000

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = ta.FocusedStyle.CursorLine.Copy().UnsetBackground()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 3)
	vp.SetContent("Type a prompt and press Enter to ask Clyde AI.")
	vp.MouseWheelEnabled = true

	vp.KeyMap = viewport.KeyMap(KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	})

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = UserStyle

	return model{
		textarea: ta,
		viewport: vp,
		spinner:  sp,
		messages: []string{},
		lastMsg:  "",
		err:      nil,
		waiting:  false,
		mouse:    true,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 10
		m.viewport.Width = msg.Width - 6
		m.textarea.SetWidth(msg.Width - 6)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.mouse {
				m.mouse = false
				return m, tea.Sequence(tea.DisableMouse, getLogCmd("Disabled mouse scroll/clicks", Info))
			} else {
				m.mouse = true
				return m, tea.Sequence(tea.EnableMouseCellMotion, getLogCmd("Enabled mouse scroll/clicks", Info))
			}
		case tea.KeyCtrlX:
			WriteClipboard(m.lastMsg)
			return m, getLogCmd(fmt.Sprintf("Copied %d characters!", len(m.lastMsg)), Info)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			prompt := m.textarea.Value()
			go AskClyde(prompt)

			m.waiting = true
			m.messages = append(m.messages, UserStyle.Render(prompt))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			return m, tea.Batch(tiCmd, vpCmd, spinner.Tick)
		}

	case errMsg:
		m.err = msg
		return m, nil

	case DiscordMessage:
		parsed := strings.ReplaceAll(msg.Content, fmt.Sprintf("<@!%s>", CurrentUserID), "`@You`")
		var md string

		m.lastMsg = parsed

		if os.Getenv("GLAMOUR_STYLE") != "" {
			md, _ = glamour.RenderWithEnvironmentConfig(parsed)
		} else {
			md, _ = glamour.Render(parsed, "dark")
		}

		m.waiting = false
		m.messages = append(m.messages, strings.Trim(md, "\n")+"\n")
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

		return m, nil

	case logMsg:
		switch msg.Type {
		case Info:
			m.messages = append(m.messages, InfoLogStyle.Render(msg.Msg))
		case Warning:
			m.messages = append(m.messages, WarningLogStyle.Render(msg.Msg))
		case Error:
			m.messages = append(m.messages, ErrorLogStyle.Render(msg.Msg))
		}

		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m model) View() string {
	var bottomView string
	if m.waiting {
		bottomView = m.spinner.View() + FadedStyle.Render(" Waiting for Clyde...")
	} else {
		bottomView = m.textarea.View()
	}

	view := ContainerStyle.Render(
		fmt.Sprintf(
			"%s\n\n%s",
			m.viewport.View(),
			lipgloss.NewStyle().Width(m.viewport.Width).Render(bottomView),
		),
	)

	return view + "\n" + HelpStyle.Render(help)
}

func RunTUI() {
	var m = initialModel()
	tui = tea.NewProgram(m, tea.WithMouseCellMotion())

	if _, err := tui.Run(); err != nil {
		panic(err)
	}
}

func getLogCmd(msg string, msgType LogType) tea.Cmd {
	return func() tea.Msg {
		return logMsg{Msg: msg, Type: msgType}
	}
}
