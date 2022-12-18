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
	configYaml, err := configuration.NewYamlConfiguration("conf.yml")
	configXml, err := configuration.NewXmlConfiguration("conf.xml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", configXml)

	fmt.Printf("%v\n", configYaml)

	processManager := process.NewManager(configYaml)
	healthCheckManager := healthcheck.NewCheckers(configYaml)

	err = processManager.StartProcesses()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5) // LAME

	var wg sync.WaitGroup

	healthCheckManager.RunChecks(processManager, &wg)
	wg.Wait()
}
