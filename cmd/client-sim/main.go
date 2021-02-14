package main

import (
	"flag"
	"fmt"
	"github.com/nspiguelman/zeus-client/pkg/client"
	"log"
	"sync"
	"time"
)

func main() {
	pin := client.ProcessCSV()

	n := flag.Int("n", 1, "number of simulated clients")
	flag.Parse()

	log.Println("Simulated clients: ", *n)
	log.Println("Room PIN: ", pin)

	var wg sync.WaitGroup
	clients := make([]*client.Client, *n)

	// instanciar a los clientes e ingresar a la sala
	seed := time.Now().Unix()
	log.Println("Initializing clients...")
	for i, _ := range clients {
		wg.Add(1)
		go func(i int) {
			username := fmt.Sprintf("user%x_%05d", seed, i)
			c, err := client.NewClient(username, pin)
			if err != nil {
				log.Fatal("construct:", err)
				return
			}
			clients[i] = c
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Println("Let the game begin!")
	// jugar
	for _, c := range clients {
		wg.Add(1)
		go c.Play(&wg)
	}
	wg.Wait()
}
