package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	token     string = "ATTAa686917cfdf4e1b56bb8635444468f4f147728183ff9c6ed4efec6070682c06eEA0B2B81"
	key       string = "34b528cb52b148c3df7f1d3c636336e1"
	trelloUrl string = "https://api.trello.com/1/"
)

type Board struct {
	Url       string `json:"url"`
	Name      string `json:"name"`
	Id        string `json:"id"`
	ShortLink string `json:"shortLink"`
}

type Card struct {
	ID           string `json:"id"`
	Descriptions string `json:"desc"`
	Name         string `json:"name"`
	ShortUrl     string `json:"shortUrl"`
	ColumnID     string `json:"idList"`
	//	ToDO: add badges
}

type Column struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FinalColumn struct {
	Column
	Cards []Card
}

type FinalBoard struct {
	Board
	Columns []FinalColumn
}

type TrelloClient struct {
	baseurl  string
	apiKey   string
	apiToken string
}

func NewTrelloClient() *TrelloClient {
	return &TrelloClient{
		apiKey:   key,
		apiToken: token,
		baseurl:  trelloUrl,
	}
}
func (t TrelloClient) getR(url string, responseStruct interface{}) (int, error) {
	query := map[string]string{
		"key":   t.apiKey,
		"token": t.apiToken,
	}
	formattedUrl := fmt.Sprintf("%s%s", t.baseurl, url)

	req, err := http.NewRequest("GET", formattedUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return http.StatusBadRequest, err
	}
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return response.StatusCode, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return response.StatusCode, nil
}

func (t TrelloClient) postR(url string, payload []byte) (int, error) {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.StatusCode, err
	}
	defer response.Body.Close()
	return response.StatusCode, nil
}

func (t TrelloClient) GetMemberBoards() ([]Board, error) {
	var boards []Board
	_, err := t.getR(fmt.Sprintf("members/me/boards"), &boards)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return boards, nil
}

func (t TrelloClient) GetBoardCards(boardShortLink string) ([]Card, error) {
	var cards []Card
	_, err := t.getR(fmt.Sprintf("boards/%s/cards/all", boardShortLink), &cards)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cards, nil

}

func (t TrelloClient) GetBoardColumns(boardShortLink string) ([]Column, error) {
	var columns []Column
	fmt.Println("ID", boardShortLink)
	_, err := t.getR(fmt.Sprintf("boards/%s/lists", boardShortLink), &columns)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return columns, nil
}
