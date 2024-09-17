package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Stage int

const (
  UsernameStage Stage = iota
  LoadingStage
  GameStage
)

type model struct {
  stage      Stage
  textInput  textinput.Model
  username   string
  width      int
  height     int
  err        error
  txtStyle   lipgloss.Style
  errorStyle lipgloss.Style
  quitStyle  lipgloss.Style
  view       string
}

type Game struct {

}
