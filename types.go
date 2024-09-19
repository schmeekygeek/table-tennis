package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Stage int
type RoomStatus int

var (
  gameServer *GameServer
  currentSeq string // current room seq
  currentRoomStat RoomStatus
)

const (
  canvasY = 20
  canvasX = 60
  LoadingStage Stage = iota
  GameStage
  Empty RoomStatus = iota
  Half
)

type model struct {
  stage      Stage
  width      int
  height     int
  err        error
  txtStyle   lipgloss.Style
  errorStyle lipgloss.Style
  quitStyle  lipgloss.Style
  view       string
  spinner    spinner.Model
  room       string // The room to which the client belongs (initially empty)
}

type Point struct {
  x, y int
}

type Room struct {
  player1    *model // player 1
  player2    *model // player 2
  ballPos    Point // position of the ball
}

type GameServer struct {
  rooms      map[string]Room // map roomId to room
}
