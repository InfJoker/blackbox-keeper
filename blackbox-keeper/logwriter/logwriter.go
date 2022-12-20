package logwriter

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/logwriter/rabbit"
	"blackbox-keeper/process"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type logWriter struct {
	stderr   io.Writer
	stdout   io.Writer
	errSleep time.Duration
}

type logWriters map[string]*logWriter

func NewLogWriters(config configuration.Config) logWriters {
	lws := make(logWriters, len(config.Apps))
	for name, appConfig := range config.Apps {
		lw := logWriter{}
		if appConfig.Exporter.Rabbit.Url != "" { // TODO Refactor this dirty if statements
			var err error
			if appConfig.Exporter.Rabbit.ErrQueue != "" {
				lw.stderr, err = rabbit.NewWriter(
					appConfig.Exporter.Rabbit.Url,
					appConfig.Exporter.Rabbit.ErrQueue,
				)
				if err != nil {
					log.Printf("Failed connecting to rabbitmq due to %s\nNot writing stderr logs for %s", err, name) // TODO it in more smart way
				}
			} else {
				log.Printf("No queue provided for %s, stderr escalated", name)
				lw.stderr = os.Stderr
			}

			if appConfig.Exporter.Rabbit.OutQueue != "" {
				lw.stdout, err = rabbit.NewWriter( // TODO we should reuse previous connection
					appConfig.Exporter.Rabbit.Url,
					appConfig.Exporter.Rabbit.OutQueue,
				)
				if err != nil {
					log.Printf("Failed connecting to rabbitmq due to %s\nNot writing stdout logs for %s", err, name) // TODO it in more smart way
				}
			} else {
				log.Printf("No queue provided for %s, stdout escalated", name)
				lw.stdout = os.Stdout
			}
		} else {
			lw.stdout = os.Stdout
			lw.stderr = os.Stderr
		}
		lw.errSleep = time.Millisecond * time.Duration(appConfig.Exporter.ErrorSleepMilli)
		lws[name] = &lw
	}
	return lws
}

func (lws logWriters) RunWriters(pm process.Manager, wg *sync.WaitGroup) {
	wg.Add(len(lws))
	for name, lw := range lws {
		go lw.run(name, pm[name], wg)
	}
}

func (lw *logWriter) run(name string, p *process.Process, wg *sync.WaitGroup) {
	defer wg.Done()
	var localWg sync.WaitGroup
	localWg.Add(2)
	go func() {
		defer localWg.Done()
		for {
			_, err := io.Copy(lw.stderr, p.GetStdErr())
			if err != nil {
				log.Printf("Failed writing stderr logs for %s due to %s", name, err)
				time.Sleep(lw.errSleep)
			}
		}
	}()

	go func() {
		defer localWg.Done()
		for {
			_, err := io.Copy(lw.stdout, p.GetStdOut())
			if err != nil {
				log.Printf("Failed writing stdout logs for %s due to %s", name, err)
				time.Sleep(lw.errSleep)
			}
		}
	}()
	localWg.Wait()
}
