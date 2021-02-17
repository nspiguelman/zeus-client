package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type InteractiveClient struct {
	username string
	pin      int
	token    string
	conn     *websocket.Conn
}

func NewInteractiveClient(pin int) (*InteractiveClient, error){
	var username string
	fmt.Printf("username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return &InteractiveClient{}, fmt.Errorf("input error: %w", err)
	}
	return &InteractiveClient{username: username, pin: pin}, nil
}

func (c *InteractiveClient) Username() string {
	return c.username
}

func (c *InteractiveClient) Pin() int {
	return c.pin
}

func (c *InteractiveClient) SetToken(token string) {
	c.token = token
}

func (c *InteractiveClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *InteractiveClient) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *InteractiveClient) Answer(question Question) Answer {
	var i int
	fmt.Printf("question id #%d: choose from 1 to %d: ", question.QuestionId, len(question.AnswerIds))
	_, err := fmt.Scanln(&i)
	if err != nil {
		fmt.Println("input error:", err, "skipping...")
		return Answer{}
	} else if i < 1 || i > len(question.AnswerIds) {
		fmt.Println("invalid input:", i, "skipping...")
		return Answer{}
	}
	ans := Answer{
		QuestionId: question.QuestionId,
		AnswerId:   question.AnswerIds[i-1],
	}
	return ans
}

func (c *InteractiveClient) PrintScore(scores map[string]Score) {
	score := scores[c.token]
	if score.IsCorrect {
		fmt.Println("username:", c.username, "- correct! score:", score.Score)
	} else {
		fmt.Println("username:", c.username, "- incorrect! score:", score.Score)
	}

}

func (c *InteractiveClient) GameOver() {
	log.Println("username:", c.username, "- closing connection...")
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
