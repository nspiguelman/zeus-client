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
	n := flag.Int("n", 1, "number of simulated clients")
	pin := flag.String("pin", "", "room pin")
	flag.Parse()

	log.Println("Simulated clients: ", *n)
	log.Println("Room PIN: ", *pin)

	clients := make([]*client.SimulatedClient, *n)

	// instanciar a los clientes e ingresar a la sala
	seed := time.Now().Unix()
	log.Println("Initializing clients...")
	for i := range clients {
		username := fmt.Sprintf("user%x_%05d", seed, i)
		clients[i] = client.NewSimulatedClient(username, *pin)
	}

	// jugar
	var wg sync.WaitGroup

	log.Println("Joining game...")
	for _, c := range clients {
		wg.Add(1)
		go func(c client.Client) {
			err := client.Login(c)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(c)
	}
	wg.Wait()


	log.Println("Starting game...")
	for _, c := range clients {
		wg.Add(1)
		go func(c client.Client) {
			client.Play(c)
			wg.Done()
		}(c)
	}
	wg.Wait()

}
