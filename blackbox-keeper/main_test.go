package main_test

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/healthcheck"
	"blackbox-keeper/process"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func startTestApp() error {
	output, err := exec.Command("go", "build", "-o", "./test/app", "./test/main.go").CombinedOutput()
	if err != nil {
		return err
	}
	log.Print(output)
	log.Print("test app successfully started")
	return nil
}

func Test(t *testing.T) {
	if err := startTestApp(); err != nil {
		t.Fatalf("failed to start test app: %v", err)
	}
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
