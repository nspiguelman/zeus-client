package main

import (
	"flag"
	"log"
	"strconv"
	"sync"
)

func main() {
	n := flag.Int("n", 1, "number of simulated clients")
	pin := flag.Int("pin", 1234, "kahoot game pin")

	flag.Parse()

	log.Println("Simulated clients: ", *n)
	log.Println("Room PIN: ", *pin)

	clients := make([]*KahootClient, *n)
	var wg sync.WaitGroup

	// instanciar a los clientes e ingresar a la sala
	log.Println("Initializing clients...")
	for i, _ := range clients {
		wg.Add(1)
		go func(i int) {
			kc, err := newKahootClient("user"+strconv.Itoa(i), *pin)
			if err != nil {
				log.Fatal("create:", err)
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
