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
	fNClient := flag.Int("n", 0, "number of simulated clients")
	fPin := flag.Int("pin", 0, "room PIN")
	flag.Parse()

	interactiveMode := *fInteractiveMode
	nClient := *fNClient
	pin := *fPin

	if *fPin == 0 {
		pin = csv.ProcessCSV()
	}

	log.Println("room PIN:", pin)
	log.Println("simulated clients:", nClient)
	log.Println("interactive mode:", interactiveMode)

	clients := make([]*client.SimulatedClient, nClient)

	// instanciar a los clientes e ingresar a la sala
	seed := time.Now().Unix()
	log.Println("initializing simulated clients...")
	for i := range clients {
		username := fmt.Sprintf("user%x_%05d", seed, i)
		clients[i] = client.NewSimulatedClient(username, pin)
	}

	var wg sync.WaitGroup
	log.Println("  joining room...")
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
	log.Println("  done.")


	for _, c := range clients {
		wg.Add(1)
		go func(c client.Client) {
			client.Play(c)
			wg.Done()
		}(c)
	}
	log.Println("ready to play.")


	if interactiveMode == true {
		log.SetOutput(ioutil.Discard)
		fmt.Println("")
		fmt.Println("/**********************************/")
		fmt.Println("/******** interactive mode ********/")
		fmt.Println("/**********************************/")
		fmt.Println("")
		ic, _ := client.NewInteractiveClient(pin)
		defer ic.GameOver()

		fmt.Println("  joining game...")
		err := client.Login(ic)
		if err != nil {
			panic(err)
		}
		fmt.Println("  done.")

		wg.Add(1)
		go func(c client.Client) {
			client.Play(c)
			wg.Done()
		}(ic)
		fmt.Println("ready to play.")
	}

	wg.Wait()
	fmt.Println("final")
	return
}
