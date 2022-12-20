package process

import (
	"blackbox-keeper/configuration"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type Process struct {
	cmd            *exec.Cmd
	stderrPipe     io.ReadCloser
	stdoutPipe     io.ReadCloser
	WaitAfterStart time.Duration // Millisecond
	RepeatAfter    time.Duration // Millisecond
	Timeout        time.Duration // Millisecond
}

type StdReader struct {
	p      *Process
	reader io.Reader
}

func (p *Process) isAlive() bool { // Note this is synchronises with pipe closing
	return p.cmd.Process != nil
}

func (r *StdReader) Read(p []byte) (int, error) {
	if !r.p.isAlive() {
		return 0, fmt.Errorf("Process was killed")
	}
	return r.reader.Read(p)
}

func (p *Process) GetStdOut() io.Reader {
	return &StdReader{
		p:      p,
		reader: p.stdoutPipe,
	}
}

func (p *Process) GetStdErr() io.Reader {
	return &StdReader{
		p:      p,
		reader: p.stderrPipe,
	}
}

func (p *Process) Start() error {
	var err error
	p.stderrPipe, err = p.cmd.StderrPipe()
	if err != nil {
		return err
	}
	p.stdoutPipe, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = p.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) Kill() error {
	p.stdoutPipe.Close() // Close pipes before lock, because reads from pipe are blocking note synchronises-with Wait
	p.stderrPipe.Close()
	err := p.cmd.Process.Kill()
	if err != nil {
		return err
	}
	p.cmd.Process.Wait()
	p.cmd.Process = nil
	p.cmd.Stderr = nil
	p.cmd.Stdout = nil
	return nil
}

type Manager map[string]*Process

func NewManager(config configuration.Config) Manager {
	res := make(Manager, len(config.Apps))
	for name, appConfig := range config.Apps {
		res[name] = &Process{
			cmd:            exec.Command(appConfig.Command),
			WaitAfterStart: time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			RepeatAfter:    time.Millisecond * time.Duration(appConfig.HealthCheck.Http.WaitAfterStartMilli),
			Timeout:        time.Millisecond * time.Duration(appConfig.HealthCheck.Http.TimeoutMilli),
		}
	}
	return res
}

func (m Manager) StartProcesses() error {
	for name, p := range m {
		err := p.Start()
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
