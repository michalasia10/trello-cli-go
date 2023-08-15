package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"trello-cli-go/handlers"
)

const APPEND = -1

type Column struct {
	focus  bool
	Status Status
	List   list.Model
	height int
	width  int
}

func (c *Column) Focus() {
	c.focus = true
}

func (c *Column) Blur() {
	c.focus = false
}

func (c *Column) Focused() bool {
	return c.focus
}

func NewColumn(status Status) Column {
	var focus bool
	if status == Todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return Column{focus: focus, Status: status, List: defaultList}
}

// Init does initial setup for the column.
func (c Column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.List.SetSize(msg.Width/Margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, handlers.Keys.Edit):
			if len(c.List.VisibleItems()) != 0 {
				task := c.List.SelectedItem().(Task)
				f := NewForm(task.Title(), task.Description())
				f.Index = c.List.Index()
				f.Col = c
				return f.Update(nil)
			}
		case key.Matches(msg, handlers.Keys.New):
			f := NewDefaultForm()
			f.Index = APPEND
			f.Col = c
			return f.Update(nil)
		case key.Matches(msg, handlers.Keys.Delete):
			return c, c.DeleteCurrent()
		case key.Matches(msg, handlers.Keys.Enter):
			return c, c.MoveToNext()
		}
	}
	c.List, cmd = c.List.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.getStyle().Render(c.List.View())
}

func (c *Column) DeleteCurrent() tea.Cmd {
	if len(c.List.VisibleItems()) > 0 {
		c.List.RemoveItem(c.List.Index())
	}

	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)
	return cmd
}

func (c *Column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		return c.List.SetItem(i, t)
	}
	return c.List.InsertItem(APPEND, t)
}

func (c *Column) setSize(width, height int) {
	c.width = width / Margin
}

func (c *Column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Task
}

func (c *Column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.List.SelectedItem().(Task); !ok {
		return nil
	}
	// move item
	c.List.RemoveItem(c.List.Index())
	task.Status = c.Status.GetNext()

	// refresh list
	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}
