package main

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheck"
	"blackbox-keeper/process"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	config, err := configuration.NewConfiguration("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", config)

	processManager, err := process.NewManager(config)
	if err != nil {
		log.Fatal(err)
	}
	healthCheckManager := healthcheck.NewCheckers(config)

	err = processManager.StartProcesses()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5) // LAME

	var wg sync.WaitGroup

	healthCheckManager.RunChecks(processManager, &wg)
	wg.Wait()
}
