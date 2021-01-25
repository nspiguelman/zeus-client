package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	n := flag.Int("n", 1, "number of simulated clients")
	pin := flag.Int("pin", 0, "kahoot game pin")
	flag.Parse()

	log.Println("Simulated clients: ", *n)
	log.Println("Room PIN: ", *pin)

	var wg sync.WaitGroup
	clients := make([]*KahootClient, *n)

	// instanciar a los clientes e ingresar a la sala
	seed := time.Now().Unix()
	log.Println("Initializing clients...")
	for i, _ := range clients {
		wg.Add(1)
		go func(i int) {
			username := fmt.Sprintf("user%x_%05d", seed, i)
			kc, err := newKahootClient(username, *pin)
			if err != nil {
				log.Fatal("construct:", err)
				return
			}
			clients[i] = kc
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Println("Let the game begin!")
	// jugar
	for _, kc := range clients {
		wg.Add(1)
		go kc.play(&wg)
	}
	wg.Wait()

	// TODO: stats

}
