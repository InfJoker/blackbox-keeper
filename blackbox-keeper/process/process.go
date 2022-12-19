package process

import (
	"blackbox-keeper/configuration"
	"blackbox-keeper/logwriter"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Process struct {
	cmd            *exec.Cmd
	Log            logwriter.LogWriter
	WaitAfterStart time.Duration // Millisecond
	RepeatAfter    time.Duration // Millisecond
	Timeout        time.Duration // Millisecond
}

func (p *Process) Start() error {
	return p.cmd.Start()
}

func (p *Process) Kill() error {
	err := p.cmd.Process.Kill()
	if err != nil {
		return err
	}
	p.cmd.Process.Wait()
	p.cmd.Process = nil
	return nil
}

type Manager map[string]*Process

func NewManager(config configuration.Config) (Manager, error) {
	res := make(Manager, len(config.Apps))
	for name, appConfig := range config.Apps {
		p := Process{
			cmd:            exec.Command(appConfig.Command),
			WaitAfterStart: time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			RepeatAfter:    time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			Timeout:        time.Millisecond * time.Duration(appConfig.HealthCheck.Http.TimeoutMilli),
		}
		stdout, err := p.cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("error on setting Stdout pipe for %s: %w", name, err)
		}
		stderr, err := p.cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("error on setting Stderr pipe for %s: %w", name, err)
		}
		p.Log = logwriter.NewLogWriter(stdout, stderr)
		res[name] = &p
	}
	return res, nil
}

func (m Manager) StartProcesses() error {
	for name, p := range m {
		err := p.cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err) // TODO Maybe abort others too here
		}
		log.Printf("%s succesfully launched!!!\n", name)
	}
	return nil
}

func (m Manager) StartProcess(name string) error {
	return m[name].Start()
}

func (m Manager) KillProcess(name string) error {
	return m[name].Kill()
}
