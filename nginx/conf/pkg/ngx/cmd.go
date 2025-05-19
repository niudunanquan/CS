package ngx

import (
	"om/pkg/util"
	"os/exec"
)

const (
	defBinary = "/usr/bin/openresty"
)

type NginxCommand struct {
	Binary string
}

func NewNginxCommand() NginxCommand {
	command := NginxCommand{
		Binary: defBinary,
	}

	return command
}

func (nc NginxCommand) ExecCommand(args ...string) *exec.Cmd {
	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, "-p", util.NginxDir)
	cmdArgs = append(cmdArgs, args...)

	return exec.Command(nc.Binary, cmdArgs...)
}

func (nc NginxCommand) Start() ([]byte, error) {
	return nc.ExecCommand().CombinedOutput()
}

func (nc NginxCommand) Test() ([]byte, error) {
	return nc.ExecCommand("-t").CombinedOutput()
}

func (nc NginxCommand) Reload() ([]byte, error) {
	return nc.ExecCommand("-s", "reload").CombinedOutput()
}

func (nc NginxCommand) Stop() ([]byte, error) {
	return nc.ExecCommand("-s", "stop").CombinedOutput()
}
