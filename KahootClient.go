package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

type KahootClient struct {
	conn     *websocket.Conn
	pin      int
	username string
	token    string
}

type KahootMessage struct {
	TypeMessage string `json:"typeMessage"`
}

type KahootScore struct {
	typeMessage  string `json:"typeMessage"`
	PartialScore int
	IsCorrect    bool
}

type KahootQuestion struct {
	typeMessage string `json:"typeMessage"`
	QuestionId  int    `json:"questionId"`
	Answerids   []int  `json:"answerIds"`
}

/*
type KahootMessage struct {
	typeMessage  string
	partialScore int
	isCorrect    bool
	questionId   int
	answerIds    []int
}
*/

type KahootAnswer struct {
	QuestionId int
	AnswerId   int
}

type KahootLogin struct {
	Token string
}

func (kc *KahootClient) login() error {
	urlStr := fmt.Sprintf("http://localhost:8080/room/%d/name/%s/login", kc.pin, kc.username)
	resp, err := http.PostForm(urlStr, url.Values{"name": {kc.username}})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var message KahootLogin
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		log.Fatal(resp.Body, err)
		// read response error
		return err
	}
	kc.token = message.Token
	return nil
}

func (kc *KahootClient) connect() error {
	urlStr := fmt.Sprintf("ws://localhost:8080/room/%d/ws", kc.pin)
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, http.Header{"token": []string{kc.token}})
	if err != nil {
		return err
	}
	kc.conn = conn
	return nil
}

func newKahootClient(username string, pin int) (*KahootClient, error) {
	var err error

	kc := KahootClient{pin: pin, username: username}

	// login
	err = kc.login()
	if err != nil {
		log.Fatal("login:", err)
		return nil, err
	}

	// establece ws
	err = kc.connect()
	if err != nil {
		log.Fatal("login:", err)
		return nil, err
	}
	log.Println("username:", kc.username, "- pin:", kc.pin, "- token:", kc.token)

	return &kc, nil
}

func (kc *KahootClient) playTrivia(questionId int, answerIds []int) {
	//time.Sleep(5000 * time.Millisecond)

	ans := KahootAnswer{
		QuestionId: questionId,
		AnswerId:   answerIds[rand.Intn(len(answerIds))],
	}
	log.Println("ans:", ans)

	b, err := json.Marshal(ans)
	if err != nil {
		log.Println("json:", err)
		return
	}

	err = kc.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("ws:", err)
		return
	}

}

func (kc *KahootClient) endGame() {
	// cuando finaliza el juego, cierra el ws
	err := kc.conn.WriteMessage(websocket.TextMessage, []byte("Player "+kc.username+": I don't wanna go Mr. Stark!"))
	if err != nil {
		log.Fatal(err)
	}
	err = kc.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (kc *KahootClient) play(wg *sync.WaitGroup) {
	defer wg.Done()
	// loop de juego
	for {
		// esperar mensaje
		// TODO: buscar una mejor manera para esperar mensaje
		_, rawMessage, err := kc.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var message KahootMessage
		err = json.Unmarshal(rawMessage, &message)
		log.Println("rawMessage:", string(rawMessage))
		if err != nil {
			log.Println(string(rawMessage))
			log.Println("unmarshall:", err)
			continue
		}

		// seleccionar tipo y responder
		switch message.TypeMessage {
		case "question":
			var kq KahootQuestion
			_ = json.Unmarshal(rawMessage, &kq)
			kc.playTrivia(kq.QuestionId, kq.Answerids)

		case "score":
			continue

		case "endgame":
			kc.endGame()
			return

		default:
			log.Printf("recv: %s", rawMessage)
		}
	}

}
