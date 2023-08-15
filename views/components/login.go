package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"trello-cli-go/handlers"
)

type (
	errMsg error
)
type Login struct {
	textInput textinput.Model
	help      help.Model
	err       error
	complete  bool
}

func InitialModelLogin() Login {
	ti := textinput.New()
	ti.Placeholder = "ApiKey"
	ti.Focus()
	ti.Width = 30

	return Login{
		textInput: ti,
		help:      help.New(),
		err:       nil,
		complete:  false,
	}
}

func (l Login) Init() tea.Cmd {
	return textinput.Blink
}

func (l Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return l, tea.Quit
		case tea.KeyCtrlQ:
			l.complete = true
			return l, cmd
		}

	// We handle errors just like any other message
	case errMsg:
		l.err = msg
		return l, nil
	}
	l.textInput, cmd = l.textInput.Update(msg)
	return l, cmd
}

func (l Login) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf(
			"What is your apikey to Trelo?\n\n%s\n\n%s",
			l.textInput.View(),
			"(esc to quit)",
		)+"\n",
		l.help.View(handlers.Keys))
}
