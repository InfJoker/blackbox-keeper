package starter

import (
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	name string
	args []string
	// add more fields for configuration
}

func NewCommand(name string, args []string) *Command {
	return &Command{name: name, args: args}
}

// Observer stores information abour running processes
type Observer struct {
	// pid -> *os.Process
	Processes map[int]*os.Process
}

// NewObserver allocates new *observer
func NewObserver() *Observer {
	return &Observer{Processes: make(map[int]*os.Process)}
}

func (o *Observer) Start(command *Command) (int, error) {
	cmd := exec.Command(command.name, command.args...)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("%s didn't start: %w", command.name, err)
	}
	o.Processes[cmd.Process.Pid] = cmd.Process
	return cmd.Process.Pid, nil
}
