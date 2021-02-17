package main

import (
	"flag"
	"fmt"
	"github.com/nspiguelman/zeus-client/pkg/client"
	"github.com/nspiguelman/zeus-client/pkg/csv"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func main() {
	fInteractiveMode := flag.Bool("i", false, "enable interactive mode")
	fNClient := flag.Int("n", 1, "number of simulated clients")
	flag.Parse()

	interactiveMode := *fInteractiveMode
	nClient := *fNClient
	pin := csv.ProcessCSV()

	log.Println("Room PIN:", pin)
	log.Println("Simulated clients:", nClient)
	log.Println("Interactive mode:", interactiveMode)

	if interactiveMode == true {
		log.SetOutput(ioutil.Discard)
	}

	clients := make([]*client.SimulatedClient, nClient)

	// instanciar a los clientes e ingresar a la sala
	seed := time.Now().Unix()
	log.Println("Initializing clients...")
	for i := range clients {
		username := fmt.Sprintf("user%x_%05d", seed, i)
		clients[i] = client.NewSimulatedClient(username, pin)
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
	if interactiveMode == true {
		iClient, _ := client.NewInteractiveClient(pin)
		err := client.Login(iClient)
		if err != nil {
			fmt.Println("could not login:", err)
			return
		}
		fmt.Println("login successful.")
		fmt.Println("starting game...")
		client.Play(iClient)
	}

	wg.Wait()
}
