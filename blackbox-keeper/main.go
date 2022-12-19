package main

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheckmanager"
	"blackbox-keeper/processmanager"
	"fmt"
	"log"
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

	processManager := processmanager.NewManager(configYaml)
	healthCheckManager := healthcheckmanager.NewManager(configYaml)

	err = processManager.StartProcesses()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5) // LAME

	for true {
		healthCheckManager.RunChecks(processManager)
	}
}
