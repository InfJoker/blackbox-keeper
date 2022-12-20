package main

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheck"
	"blackbox-keeper/logwriter"
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

	processManager := process.NewManager(config)
	logWriters := logwriter.NewLogWriters(config)
	healthCheckManager := healthcheck.NewCheckers(config)

	err = processManager.StartProcesses()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5) // LAME

	var wg sync.WaitGroup
	logWriters.RunWriters(processManager, &wg)

	healthCheckManager.RunChecks(processManager, &wg)

	wg.Wait()
}
