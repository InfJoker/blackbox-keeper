package processmanager

import (
	"blackbox-keeper/configuration"
	"log"
	"os"
	"os/exec"
)

type ProcessManager struct {
	commands  map[string]*exec.Cmd
	processes map[string]*os.Process
}

func NewManager(config configuration.ConfigYaml) *ProcessManager {
	commands := map[string]*exec.Cmd{}
	for name, appConfig := range config.Apps {
		commands[name] = exec.Command(appConfig.Command) // TODO Split command into args
	}
	return &ProcessManager{commands: commands, processes: make(map[string]*os.Process)}
}

func (p *ProcessManager) StartProcesses() error {
	for name, command := range p.commands {
		err := command.Start()
		if err != nil {
			return err // TODO Maybe abort others too here
		}
		log.Printf("%s succesfully launched!!!\n", name)
		p.processes[name] = command.Process
	}
	return nil
}

func (p *ProcessManager) StartProcess(name string) error {
	return p.commands[name].Start()
}

func (p *ProcessManager) KillProcess(name string) error {
	err := p.processes[name].Kill()
	if err != nil {
		return err
	}
	p.processes[name].Wait()
	p.commands[name].Process = nil
	return nil
}
