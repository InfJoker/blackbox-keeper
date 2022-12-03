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
    config, err := configuration.NewConfiguration("conf.yml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%v/n", config)

    processManager := processmanager.NewManager(config)
    healthCheckManager := healthcheckmanager.NewManager(config)

    err = processManager.StartProcesses()
    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(time.Second * 5) // LAME

    for true {
        healthCheckManager.RunChecks(processManager)
    }
}
