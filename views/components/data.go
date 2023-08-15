package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"trello-cli-go/api"
	"trello-cli-go/handlers"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type modelList struct {
	list list.Model
}

func (m modelList) Init() tea.Cmd {
	return nil
}

func (m modelList) Update(msg tea.Msg) (modelList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	default:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelList) View() string {
	return docStyle.Render(m.list.View())
}

type StatusResp int
type DataBoard struct {
	spinnerModel  spinner.Model
	help          help.Model
	err           error
	fetchComplete bool
	fetching      bool
	FinalBoard    api.FinalBoard
	BoardsToPick  []api.Board
	pickedBoard   api.Board
	apiClient     api.TrelloClient
	showSpinner   bool
	options       modelList
}

func InitialDataBoard(configuredApiClient *api.TrelloClient) DataBoard {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return DataBoard{
		spinnerModel:  spin,
		help:          help.New(),
		err:           nil,
		fetchComplete: false,
		showSpinner:   true,
		fetching:      false,
		apiClient:     *configuredApiClient,
	}
}

func (d *DataBoard) Init() tea.Cmd {
	boards, err := d.apiClient.GetMemberBoards()
	if err != nil {
		return d.spinnerModel.Tick
	}

	d.BoardsToPick = boards
	var listRows []list.Item
	for _, board := range d.BoardsToPick {
		listRows = append(listRows, item{title: board.Name, id: board.Id})
	}

	d.options = modelList{list: list.New(listRows, list.NewDefaultDelegate(), 100, 10)}
	d.options.list.Title = "Boards"
	return d.spinnerModel.Tick
}

func (d *DataBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return d, tea.Quit
		case tea.KeyEnter:
			d.fetchComplete = true
			selectedItem := d.options.list.SelectedItem()
			if sI, ok := selectedItem.(item); ok {
				for _, board := range d.BoardsToPick {
					if board.Id == sI.id {
						d.pickedBoard = board
						b := handlers.NewBoardDataBuilder(d.apiClient)
						d.FinalBoard, _ = b.BuildBoardData(board.ShortLink)
						d.fetchComplete = true
						return d, tea.Quit
					}
				}
			}

			return d, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		if len(d.options.list.Items()) > 0 {
			d.options.list.SetSize(msg.Width-h, msg.Height-v)
		}
		// We handle errors just like any other message
	case errMsg:
		d.err = msg
		return d, nil
	default:
		d.options, cmd = d.options.Update(msg)
		return d, cmd
	}

	d.options, cmd = d.options.Update(msg)
	return d, cmd
}

func (d *DataBoard) View() string {
	if len(d.BoardsToPick) <= 0 {
		return fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", d.spinnerModel.View())
	}
	return d.options.View()
}
