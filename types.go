package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type Stage int

var gameServer *GameServer

const (
  LoadingStage Stage = iota
  GameStage
)

type model struct {
  stage      Stage
  username   string
  width      int
  height     int
  err        error
  txtStyle   lipgloss.Style
  errorStyle lipgloss.Style
  quitStyle  lipgloss.Style
  view       string
  spinner    spinner.Model
}

type Client struct {
  session   *ssh.Session // The session of the client
  pos       Point // The position of the player
  room      string // The room to which the client belongs (initially empty)
}

type Point struct {
  x, y int
}

type Room struct {
  player1    Client // player 1
  player2    Client // player 2
  ballPos    Point // position of the ball
}

type GameServer struct {
  clients    []*ssh.Session  // list of all connected (but not matched) clients
  rooms      map[string]Room // map roomId to room
}
