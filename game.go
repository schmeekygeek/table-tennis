package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func InitGameServer() {
  gameServer = &GameServer{
  	rooms: map[string]Room{},
  }
  currentRoomStat = Empty
  currentSeq = randSeq(5)
}

func randSeq(n int) string {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[r.Intn(len(letters))]
  }
  return string(b)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
  pty, _, _ := s.Pty()

  renderer := bubbletea.MakeRenderer(s)
  quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
  errorStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))
  txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))

  // spinner
  sp := spinner.New()
	sp.Spinner = spinner.Globe
	sp.Style = errorStyle

  m := model{
  	stage:      LoadingStage,
  	width:      pty.Window.Width,
  	height:     pty.Window.Height,
  	err:        nil,
  	txtStyle:   txtStyle,
  	errorStyle: errorStyle,
  	quitStyle:  quitStyle,
  	view:       "",
  	spinner:    sp,
  }

  switch currentRoomStat {
  case Empty:
    room := Room{
    	player1: &model{},
    	player2: &model{},
    	ballPos: Point{},
    }
    room.player1 = &m
    currentRoomStat = Half
    gameServer.rooms[currentSeq] = room
  case Half:
    if room, ok := gameServer.rooms[currentSeq]; ok {
      room.player2 = &m
      currentRoomStat = Empty
      gameServer.rooms[currentSeq] = room
      currentSeq = randSeq(5)
      // next step: move to game stage
    }
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
      fmt.Println(gameServer.rooms)
      fmt.Println(currentRoomStat)
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
