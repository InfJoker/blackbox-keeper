package main

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheck"
	"blackbox-keeper/process"
	"log"
	"sync"
	"time"
)

func main() {
	config, err := configuration.NewConfiguration("conf.yml")
	_, err = configuration.NewXmlConfiguration("conf.xml")
	if err != nil {
		log.Fatal(err)
	}

	processManager := process.NewManager(config)
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
