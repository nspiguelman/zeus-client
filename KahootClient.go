package main

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
)

type KahootClient struct {
	conn     *websocket.Conn
	pin      int
	username string
	token    string
}

type KahootMessage struct {
	TypeMessage string
}

type KahootQuestion struct {
	TypeMessage string
	QuestionId  int
	AnswerIds   []int
}

type KahootAnswer struct {
	QuestionId int `json:"questionId,omitempty"`
	AnswerId   int `json:"answerId,omitempty"`
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
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, http.Header{"token": []string{kc.token}, "pin": []string{strconv.Itoa(kc.pin)}})
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

	b, err := json.Marshal(ans)
	if err != nil {
		log.Println("json:", err)
		return
	}

	err = kc.conn.WriteMessage(websocket.TextMessage, b)
	log.Println("username:", kc.username, "- answer:", string(b))
	if err != nil {
		log.Println("ws:", err)
		return
	}

}

func (kc *KahootClient) endGame() {
	log.Println("username:", kc.username, "- closing connection...")
	err := kc.conn.Close()
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
		if err != nil {
			log.Println("unmarshall:", err)
			continue
		}

		// seleccionar tipo y responder
		switch message.TypeMessage {
		case "question":
			var kq KahootQuestion
			_ = json.Unmarshal(rawMessage, &kq)
			kc.playTrivia(kq.QuestionId, kq.AnswerIds)

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
