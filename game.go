package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

func InitGameServer() {
  gameServer = new(GameServer)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
  pty, _, _ := s.Pty()
  fmt.Println(s.User())
  renderer := bubbletea.MakeRenderer(s)
  errorStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))
  txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
  quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

  // spinner
  sp := spinner.New()
	sp.Spinner = spinner.Globe
	sp.Style = errorStyle

  m := model{
  	username:   "",
  	width:      pty.Window.Width,
  	height:     pty.Window.Height,
  	err:        nil,
  	txtStyle:   txtStyle,
  	errorStyle: errorStyle,
  	quitStyle:  quitStyle,
  	view:       "",
  	spinner:    sp,
  }

  return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) View() string {
  var sb strings.Builder
  if m.stage == LoadingStage {
      sb.WriteString(fmt.Sprintf("\n %s You're being matched with another player, hang tight!", m.spinner.View()))
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
        // proceed to loading screen
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
  m.spinner, cmd = m.spinner.Update(msg)
  return m, cmd
}
