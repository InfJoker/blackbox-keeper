package process

import (
	"blackbox-keeper/configuration"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	cmd            *exec.Cmd
	Process        *os.Process
	WaitAfterStart time.Duration // Millisecond
	RepeatAfter    time.Duration // Millisecond
	Timeout        time.Duration // Millisecond
}

func (p *Process) Start() error {
	return p.cmd.Start()
}

func (p *Process) Kill() error {
	err := p.Process.Kill()
	if err != nil {
		return err
	}
	p.Process.Wait()
	p.cmd.Process = nil
	return nil
}

type Manager map[string]*Process

func NewManager(config configuration.Config) Manager {
	res := make(Manager, len(config.Apps))
	for name, appConfig := range config.Apps {
		res[name] = &Process{
			cmd:            exec.Command(appConfig.Command),
			Process:        nil,
			WaitAfterStart: time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			RepeatAfter:    time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			Timeout:        time.Millisecond * time.Duration(appConfig.HealthCheck.Http.TimeoutMilli),
		}
	}
	return res
}

func (m Manager) StartProcesses() error {
	for name, p := range m {
		err := p.cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err) // TODO Maybe abort others too here
		}
		log.Printf("%s succesfully launched!!!\n", name)
		m[name].Process = p.cmd.Process
	}
	return nil
}

func (m Manager) StartProcess(name string) error {
	return m[name].Start()
}

func (m Manager) KillProcess(name string) error {
	return m[name].Kill()
}
