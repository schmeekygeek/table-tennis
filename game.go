package main

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
  pty, _, _ := s.Pty()

  renderer := bubbletea.MakeRenderer(s)
  errorStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))
  txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
  quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

  ti := textinput.New()
  ti.Placeholder = "Pikachu"
  ti.Focus()
  ti.CharLimit = 156
  ti.Width = 20

  m := model{
  	stage:      UsernameStage,
  	textInput:  ti,
  	username:   "",
  	width:      pty.Window.Width,
  	height:     pty.Window.Height,
  	err:        nil,
  	txtStyle:   txtStyle,
  	errorStyle: errorStyle,
  	quitStyle:  quitStyle,
  	view:       "",
  }

  return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) View() string {
  var sb strings.Builder
  if m.stage == UsernameStage {
    sb.WriteString(m.txtStyle.Render("Please enter your username:"))
    sb.WriteString("\n\n")
    sb.WriteString(m.textInput.View())
    sb.WriteString("\n")
    if regexp.MustCompile(`\s`).Match([]byte(m.textInput.Value())) {
      sb.WriteString(m.errorStyle.Render("Username cannot contain whitespaces"))
      sb.WriteString("\n")
    }
    sb.WriteString("\n")
    sb.WriteString(m.quitStyle.Render("Press `esc` to quit."))
  }
  return sb.String()
}

type (
  errMsg error
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyEnter:
      if !regexp.MustCompile(`\s`).Match([]byte(m.textInput.Value())) {
        // proceed to loading screen
      }
    case tea.KeyCtrlC, tea.KeyEsc:
      return m, tea.Quit
  }
  case tea.WindowSizeMsg:
    m.height = msg.Height
    m.width = msg.Width
  case errMsg:
    m.err = msg
    return m, nil
  }

  m.textInput, cmd = m.textInput.Update(msg)

  return m, cmd
}
