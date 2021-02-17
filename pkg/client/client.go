package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	Username() string
	Pin() int
	Token() string
	SetToken(string)
	Conn() *websocket.Conn
	SetConn(*websocket.Conn)
	Answer(Question) Answer
	PrintScore(map[string]Score)
	EndGame()
}


type Config struct {
	ServerHost string
	ServerPort int
}

// TODO: usar envvar
var cfg = Config{ServerHost: "localhost", ServerPort: 8080}

func Login(client Client) error {
	pin := client.Pin()
	username := client.Username()
	urlStr := fmt.Sprintf("http://%s:%d/room/%d/name/%s/login", cfg.ServerHost, cfg.ServerPort, pin, username)

	resp, err := http.PostForm(urlStr, url.Values{"name": {username}})
	if err != nil {
		return fmt.Errorf("can't login: %w", err)
	}
	defer resp.Body.Close()

	var message map[string]interface{}
	var ws *websocket.Conn

	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return fmt.Errorf("can't decode login response: %s (%w)", resp.Body, err)
	}
	if message["token"] == nil {
		return fmt.Errorf("login response: %s", message)
	}
	token := message["token"].(string)
	urlStr = fmt.Sprintf("ws://%s:%d/room/%d/ws", cfg.ServerHost, cfg.ServerPort, pin)

	ws, _, err = websocket.DefaultDialer.Dial(urlStr,
		http.Header{"token": []string{token}, "pin": []string{strconv.Itoa(pin)}})
	if err != nil {
		return fmt.Errorf("can't establish ws: %w", err)
	}

	client.SetToken(token)
	client.SetConn(ws)
	return nil
}

func Play(client Client) {

	// loop de juego
	for {
		message, err := readMessage(client)
		if err != nil {
			log.Println("json:", err)
			return
		}

		// seleccionar tipo y responder
		switch message.TypeMessage {
		case "question":
			answer := client.Answer(
				Question{
					message.QuestionId,
					message.AnswerIds,
				})

			bAnswer, err := json.Marshal(answer)
			if err != nil {
				log.Println("can't marshal answer json:", err, "skipping...")
				continue
			}
			log.Println("username:", client.Username(), "- answer:", string(bAnswer))
			writeMessage(client, bAnswer)

		case "score":
			client.PrintScore(message.Scores)

		case "endgame":
			client.EndGame()
			return

		default:
			log.Printf("recv: %v", message)
		}
	}

}

func readMessage(client Client) (Message, error){
	// esperar mensaje
	_, rawMessage, err := client.Conn().ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return Message{}, err
	}

	var message Message
	//log.Println(string(rawMessage))
	err = json.Unmarshal(rawMessage, &message)
	if err != nil {
		log.Println("unmarshall:", err)
		return Message{}, err
	}
	return message, err
}

func writeMessage(client Client, message []byte) {

	err := client.Conn().WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("ws:", err)
		return
	}

}