package healthcheck

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/process"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Checker interface {
	Check(timeout time.Duration) error
}

type Checkers map[string]Checker

type httpChecker struct { // TODO should be with interface and in separate file
	host string
	port int64
	path string
}

var (
	NotOk   = errors.New("request to app failed with bad response")
	TimeOut = errors.New("timeout error")
)

func (h *httpChecker) Check(timeout time.Duration) error {
	requestURL := fmt.Sprintf("http://%s:%d/%s", h.host, h.port, h.path)
	r := make(chan *http.Response, 1)
	e := make(chan error, 1)
	go func() {
		res, err := http.Get(requestURL)
		if err != nil {
			e <- err
		} else {
			r <- res
		}
	}()

	select {
	case err := <-e:
		fmt.Printf("error making http request: %s\n", err)
		return err
	case res := <-r:
		if res.StatusCode != 200 {
			return NotOk
		}
	case <-time.After(timeout):
		return TimeOut
	}
	return nil
}

func NewCheckers(config configuration.Config) Checkers { // TODO parse more parameters from config especially stop-action
	checkers := make(Checkers, len(config.Apps))
	for name, appConfig := range config.Apps {
		checkers[name] = &httpChecker{
			host: appConfig.HealthCheck.Http.Host,
			port: appConfig.HealthCheck.Http.Port,
			path: appConfig.HealthCheck.Http.Path,
		}
	}
	return checkers
}

func (h Checkers) RunChecks(pm process.Manager, wg *sync.WaitGroup) <-chan struct{} {
	wg.Add(len(h))
	for name, check := range h {
		name, check := name, check
		go func() {
			defer wg.Done()
			runCheck(name, check, pm[name])
		}()
	}
	return make(chan struct{})
}

func runCheck(name string, check Checker, pm *process.Process) {
	for {
		err := check.Check(pm.Timeout)
		if err != nil { // TODO do this smarter
			log.Print(err)
			log.Printf("Restarting app %s...", name)
			err = pm.Kill()
			if err != nil {
				log.Printf("Failed stopping the app %s", name)
			}
			err = pm.Start()
			if err != nil {
				log.Printf("Failed starting the app %s", name)
			}
			log.Printf("Successfuly restarted app %s", name)
			time.Sleep(pm.WaitAfterStart) // This is lame
		} else {
			c := make(chan struct{})
			go func() {
				var wg sync.WaitGroup
				wg.Add(2)
				go func() {
					_, err := io.Copy(os.Stderr, pm.Stderr)
					if err != nil {
						log.Printf("error on writing from stderr pipe of %s: %v", name, err)
					}
					wg.Done()
				}()
				go func() {
					_, err := io.Copy(os.Stdout, pm.Stdout)
					if err != nil {
						log.Printf("error on writing from stdout pipe of %s: %v", name, err)
					}
					wg.Done()
				}()
				wg.Wait()
				c <- struct{}{}
			}()
			select {
			case <-c:
				log.Printf("log successfuly saved")
			case <-time.After(pm.RepeatAfter):
				// TODO: but it will be done anyway in background, cause idk how to kill goroutine
				log.Printf("didn't have enough time for saving logs")
			}
		}
	}
}
