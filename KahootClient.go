package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KahootClient struct {
	conn   *websocket.Conn
	pin    int
	userId int
}

func newKahootClient(userId int, pin int) *KahootClient {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:5000/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	kc := KahootClient{conn, pin, userId}
	//kc.conn.WriteMessage(websocket.TextMessage, []byte("Player "+strconv.Itoa(kc.userId)+" entered the room."))
	return &kc
}

func (kc *KahootClient) playTrivia(n int, timeout int) {
	// este es un metodo generico para trivias
	// elige un entero entre 0 y n-1 y lo envia por el ws luego de una sleep aleatorio de timeout seg
	ans := strconv.Itoa(rand.Intn(n))
	time.Sleep(time.Duration(rand.Intn(timeout*1000)) * time.Millisecond)
	kc.conn.WriteMessage(websocket.TextMessage, []byte("Player "+strconv.Itoa(kc.userId)+": "+ans))
}

func (kc *KahootClient) playMultipleChoice(timeout int) {
	// multiple choice son 4 opciones
	kc.playTrivia(4, timeout)
}

func (kc *KahootClient) playTrueOrFalse(timeout int) {
	// true = 0, false = 1
	kc.playTrivia(2, timeout)
}

func (kc *KahootClient) endGame() {
	// cuando finaliza el juego, cierra el ws
	kc.conn.WriteMessage(websocket.TextMessage, []byte("Player "+strconv.Itoa(kc.userId)+": I don't wanna go Mr. Stark!"))
	kc.conn.Close()
}

func (kc *KahootClient) play(wg *sync.WaitGroup) {
	defer wg.Done()
	// loop de juego
	for {
		// esperar broadcast
		// TODO: buscar una mejor manera para esperar broadcast
		_, message, err := kc.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		// TODO: buscar en los headers los parametros de la pregunta
		// seleccionar tipo y responder
		switch {
		case strings.Contains(string(message), "multiple_choice"):
			kc.playMultipleChoice(5)

		case strings.Contains(string(message), "true_or_false"):
			kc.playTrueOrFalse(5)

		case strings.Contains(string(message), "endgame"):
			kc.endGame()
			return

		default:
			//log.Printf("recv: %s", message)
		}
	}

	log.Println("bye!")

	//

	return
}
