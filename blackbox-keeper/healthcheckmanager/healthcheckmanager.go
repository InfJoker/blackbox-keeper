package healthcheckmanager

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/processmanager"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HealthCheckManager struct {
	checkers map[string]*httpChecker // TODO interface
}

type httpChecker struct { // TODO should be with interface and in separate file
	host string
	port int64
	path string
}

var (
	NotOk = errors.New("Request to app failed with bad response")
)

func (h *httpChecker) Check() error {
	requestURL := fmt.Sprintf("http://%s:%d/%s", h.host, h.port, h.path)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return err
	}
	if res.StatusCode != 200 {
		return NotOk
	}
	return nil
}

func NewManager(config configuration.ConfigYaml) *HealthCheckManager { // TODO parse more parameters from config especially stop-action
	checkers := map[string]*httpChecker{}
	for name, appConfig := range config.Apps {
		checkers[name] = &httpChecker{
			host: appConfig.HealthCheck.Http.Host,
			port: appConfig.HealthCheck.Http.Port,
			path: appConfig.HealthCheck.Http.Path,
		}
	}
	return &HealthCheckManager{checkers: checkers}
}

func (h *HealthCheckManager) RunChecks(proccessManager *processmanager.ProcessManager) {
	for name, check := range h.checkers {
		err := check.Check()
		if err != nil { // TODO do this smarter
			log.Printf("Restarting app %s...", name)
			err = proccessManager.KillProcess(name)
			if err != nil {
				log.Printf("Failed stopping the app %s", name)
			}
			err = proccessManager.StartProcess(name)
			if err != nil {
				log.Printf("Failed starting the app %s", name)
			}
			time.Sleep(time.Second * 5) // This is lame
			log.Printf("Successfuly restarted app %s", name)
		}
	}
}
