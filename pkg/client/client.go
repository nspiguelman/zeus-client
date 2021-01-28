package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	conn     *websocket.Conn
	pin      int
	username string
	token    string
}

type Login struct {
	Token string
}

type Config struct {
	ServerHost string
	ServerPort int
}

// TODO: usar envvar
var cfg = Config{ServerHost: "localhost", ServerPort: 8080}

func (c *Client) login() error {
	urlStr := fmt.Sprintf("http://%s:%d/room/%d/name/%s/login", cfg.ServerHost, cfg.ServerPort,c.pin, c.username)
	resp, err := http.PostForm(urlStr, url.Values{"name": {c.username}})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var message Login
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		log.Fatal(resp.Body, err)
		// read response error
		return err
	}
	c.token = message.Token
	return nil
}

func (c *Client) connect() error {
	urlStr := fmt.Sprintf("ws://%s:%d/room/%d/ws", cfg.ServerHost, cfg.ServerPort, c.pin)
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, http.Header{"token": []string{c.token}, "pin": []string{strconv.Itoa(c.pin)}})
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func NewClient(username string, pin int) (*Client, error) {
	var err error

	c := Client{pin: pin, username: username}

	// login
	err = c.login()
	if err != nil {
		log.Fatal("login:", err)
		return nil, err
	}

	// establece ws
	err = c.connect()
	if err != nil {
		log.Fatal("connect:", err)
		return nil, err
	}
	log.Println("username:", c.username, "- pin:", c.pin, "- token:", c.token)

	return &c, nil
}

func (c *Client) playTrivia(questionId int, answerIds []int) {
	//time.Sleep(5000 * time.Millisecond)

	ans := Answer{
		QuestionId: questionId,
		AnswerId:   answerIds[rand.Intn(len(answerIds))],
	}

	b, err := json.Marshal(ans)
	if err != nil {
		log.Println("json:", err)
		return
	}
	// TODO: mejorar mecanismo para determinar el tiempo
	time.Sleep(time.Duration(rand.Intn(6000))*time.Millisecond)

	err = c.conn.WriteMessage(websocket.TextMessage, b)
	log.Println("username:", c.username, "- answer:", string(b))
	if err != nil {
		log.Println("ws:", err)
		return
	}

}

func (c *Client) endGame() {
	log.Println("username:", c.username, "- closing connection...")
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) Play(wg *sync.WaitGroup) {
	defer wg.Done()
	// loop de juego
	for {
		// esperar mensaje
		// TODO: buscar una mejor manera para esperar mensaje
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var message Message
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println("unmarshall:", err)
			continue
		}

		// seleccionar tipo y responder
		switch message.TypeMessage {
		case "question":
			var q Question
			_ = json.Unmarshal(rawMessage, &q)
			c.playTrivia(q.QuestionId, q.AnswerIds)

		case "score":
			continue

		case "endgame":
			c.endGame()
			return

		default:
			log.Printf("recv: %s", rawMessage)
		}
	}

}
