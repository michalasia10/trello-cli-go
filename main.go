package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"trello-cli-go/api"
	"trello-cli-go/views/components"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	apiClient := api.NewTrelloClient()
	components.Bboard = components.NewBoard(apiClient)
	p := tea.NewProgram(components.Bboard)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//boardBuilder := handlers.NewBoardDataBuilder(*apiClient)
	//boards, err := apiClient.GetMemberBoards()
	//if err != nil {
	//	fmt.Println("error")
	//	return
	//}
	//board, err := boardBuilder.BuildBoardData(boards[1].ShortLink)
	//if err != nil {
	//	fmt.Println("ERRR")
	//}
	//fmt.Println(board)

}
