package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type InteractiveClient struct {
	username string
	pin      int
	conn     *websocket.Conn
	token    string
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
	}
	ans := Answer{
		QuestionId: question.QuestionId,
		AnswerId:   question.AnswerIds[i],
	}
	return ans
}

func (c *InteractiveClient) PrintScore() {

}

func (c *InteractiveClient) EndGame() {
	log.Println("username:", c.username, "- closing connection...")
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
