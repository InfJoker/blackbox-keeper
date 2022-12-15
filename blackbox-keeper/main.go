package main

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheck"
	"blackbox-keeper/process"
	"fmt"
	"log"
	"time"
)

func main() {
	config, err := configuration.NewConfiguration("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v/n", config)

	processManager := process.NewManager(config)
	healthCheckManager := healthcheck.NewCheckers(config)

	err = processManager.StartProcesses()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5) // LAME

	<-healthCheckManager.RunChecks(processManager)
}
