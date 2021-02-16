package client

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

type SimulatedClient struct {
	username string
	pin      int
	conn     *websocket.Conn
	token    string
}

func NewSimulatedClient(username string, pin int) *SimulatedClient {
	return &SimulatedClient{username: username, pin: pin}
}

func (c *SimulatedClient) Username() string {
	return c.username
}

func (c *SimulatedClient) Pin() int {
	return c.pin
}

func (c *SimulatedClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *SimulatedClient) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *SimulatedClient) Answer(question Question) Answer {
	time.Sleep(5000 * time.Millisecond)
	ans := Answer{
		QuestionId: question.QuestionId,
		AnswerId:   question.AnswerIds[rand.Intn(len(question.AnswerIds))],
	}
	return ans
}

func (c *SimulatedClient) PrintScore() {

}

func (c *SimulatedClient) EndGame() {
	log.Println("username:", c.username, "- closing connection...")
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
