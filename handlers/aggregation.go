package handlers

import (
	"fmt"
	"trello-cli-go/api"
)

type BoardDataBuilder struct {
	trelloClient api.TrelloClient
}

func NewBoardDataBuilder(client api.TrelloClient) *BoardDataBuilder {
	return &BoardDataBuilder{
		trelloClient: client,
	}
}
func (b *BoardDataBuilder) AggregateCardsToColumns(cards []api.Card, columns []api.Column) []api.FinalColumn {
	columnsMap := make(map[string]string) // map ColumnID to ColumnName
	for _, column := range columns {
		columnsMap[column.ID] = column.Name
	}

	columnsForCards := make(map[string][]api.Card)
	for _, card := range cards {
		if _, exists := columnsForCards[card.ColumnID]; !exists {
			columnsForCards[card.ColumnID] = make([]api.Card, 0)
		}
		columnsForCards[card.ColumnID] = append(columnsForCards[card.ColumnID], card)
	}

	var finalColumns []api.FinalColumn
	for columnID, columnCards := range columnsForCards {
		finalColumns = append(finalColumns, api.FinalColumn{
			Column: api.Column{
				ID:   columnID,
				Name: columnsMap[columnID], // Populate the Name based on ColumnID
			},
			Cards: columnCards,
		})
	}

	return finalColumns
}
func (b *BoardDataBuilder) BuildBoardData(boardShortLink string) (api.FinalBoard, error) {
	boards, err := b.trelloClient.GetMemberBoards()
	if err != nil {
		return api.FinalBoard{}, err
	}

	// Find the relevant board using boardShortLink or any other criteria
	var selectedBoard api.Board
	for _, board := range boards {
		if board.ShortLink == boardShortLink {
			selectedBoard = board
			break
		}
	}

	if selectedBoard.Id == "" {
		return api.FinalBoard{}, fmt.Errorf("board not found")
	}

	columns, err := b.trelloClient.GetBoardColumns(selectedBoard.ShortLink)
	if err != nil {
		return api.FinalBoard{}, err
	}

	cards, err := b.trelloClient.GetBoardCards(selectedBoard.ShortLink)
	if err != nil {
		return api.FinalBoard{}, err
	}

	FinalColumns := b.AggregateCardsToColumns(cards, columns)

	finalBoard := api.FinalBoard{
		Board:   selectedBoard,
		Columns: FinalColumns,
	}

	return finalBoard, nil
}
