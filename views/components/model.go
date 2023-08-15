package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"trello-cli-go/api"
	"trello-cli-go/handlers"
)

type Board struct {
	help      help.Model
	loaded    bool
	focused   Status
	cols      []Column
	quitting  bool
	login     Login
	data      DataBoard
	apiClient api.TrelloClient
}

func NewBoard(configuredApiClient *api.TrelloClient) *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: Todo, apiClient: *configuredApiClient}
}

func (m *Board) Init() tea.Cmd {
	m.login = InitialModelLogin()
	m.data = InitialDataBoard(&m.apiClient)
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//if !m.login.complete {
	//	m.login.Init()
	//	loginModel, _ := m.login.Update(msg)
	//	m.login = loginModel.(Login)
	//}

	if !m.data.fetchComplete {
		if m.login.textInput.Value() == msg {
			msg = ""
		}
		if len(m.data.BoardsToPick) <= 0 {
			m.data.Init()
		}
		_, _ = m.data.Update(msg)

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - Margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(Column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case Form:
		return m, m.cols[m.focused].Set(msg.Index, msg.CreateTask())
	case moveMsg:
		return m, m.cols[m.focused.GetNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, handlers.Keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, handlers.Keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.GetPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, handlers.Keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.GetNext()
			m.cols[m.focused].Focus()
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(Column); ok {
		m.cols[m.focused] = res.(Column)
	} else {
		return res, cmd
	}
	return m, cmd
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	//if !m.login.complete {
	//	return m.login.View()
	//}
	if !m.data.fetchComplete {
		return m.data.View()
	}
	if m.quitting {
		return ""
	}
	if !m.loaded {
		m.login = InitialModelLogin()
		return m.login.View()
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[Todo].View(),
		m.cols[InProgress].View(),
		m.cols[Done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(handlers.Keys))
}

func (m *Board) InitLists() {
	m.cols = []Column{
		NewColumn(Todo),
		NewColumn(InProgress),
		NewColumn(Done),
	}
	// Init To Do
	m.cols[Todo].List.Title = "To Do"
	m.cols[Todo].List.SetItems([]list.Item{
		Task{Status: Todo, title: "buy milk", description: "strawberry milk"},
		Task{Status: Todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{Status: Todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	m.cols[InProgress].List.Title = "In Progress"
	m.cols[InProgress].List.SetItems([]list.Item{
		NewTask(InProgress, "write code", "don't worry, it's Go"),
	})
	// Init done
	m.cols[Done].List.Title = "Done"
	m.cols[Done].List.SetItems([]list.Item{
		NewTask(Done, "stay cool", "as a cucumber"),
	})
}
