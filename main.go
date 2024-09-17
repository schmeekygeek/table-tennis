package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
  host = "localhost"
  port = "23234"
)

func main() {
  s, err := wish.NewServer(
    wish.WithAddress(net.JoinHostPort(host, port)),
    wish.WithHostKeyPath(".ssh/id_ed25519"),
    wish.WithMiddleware(
      bubbletea.Middleware(teaHandler),
      activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
      logging.Middleware(),
    ),
  )
  if err != nil {
    log.Error("Could not start server", "error", err)
  }

  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  log.Info("Starting SSH server", "host", host, "port", port)
  go func() {
    if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
      log.Error("Could not start server", "error", err)
      done <- nil
    }
    }()

  <-done
  log.Info("Stopping SSH server")
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer func() { cancel() }()
  if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
    log.Error("Could not stop server", "error", err)
  }
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
// Just a generic tea.Model to demo terminal information of ssh.

func (m model) Init() tea.Cmd {
  m.view += m.txtStyle.Render(fmt.Sprintf("Please type in a username \n\n%s", m.textInput.View()))
  m.view += "\n\n" + m.quitStyle.Render("Press 'esc' to quit\n")
  return textinput.Blink
}
