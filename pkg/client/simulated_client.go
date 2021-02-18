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
	token    string
	conn     *websocket.Conn
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

func (c *SimulatedClient) SetToken(token string) {
	c.token = token
}

func (c *SimulatedClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *SimulatedClient) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *SimulatedClient) Answer(question Question) Answer {
	time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	ans := Answer{
		QuestionId: question.QuestionId,
		AnswerId:   question.AnswerIds[rand.Intn(len(question.AnswerIds))],
	}
	return ans
}

func (c *SimulatedClient) PrintScore(scores map[string]Score) {
	score := scores[c.token]
	//log.Println(c.token, scores)
	if score.IsCorrect {
		log.Println("username:", c.username, "- correct! score:", score.Score)
	} else {
		log.Println("username:", c.username, "- incorrect! score:", score.Score)
	}
}

func (c *SimulatedClient) NextRound() {

}

func (c *SimulatedClient) GameOver() {
	log.Println("username:", c.username, "- closing connection...")
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
