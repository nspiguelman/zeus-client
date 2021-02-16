package client

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
)

type InteractiveClient struct {
	username string
	pin      string
	conn     *websocket.Conn
	token    string
}

func NewInteractiveClient(pin string) *InteractiveClient {
	// TODO: pedir username
	username := "user"
	return &InteractiveClient{username: username, pin: pin}
}

func (c *InteractiveClient) Username() string {
	return c.username
}

func (c *InteractiveClient) Pin() string {
	return c.pin
}

func (c *InteractiveClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *InteractiveClient) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *InteractiveClient) Answer(question Question) Answer {
	// TODO: Interactive answer
	ans := Answer{
		QuestionId: question.QuestionId,
		AnswerId:   question.AnswerIds[rand.Intn(len(question.AnswerIds))],
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
