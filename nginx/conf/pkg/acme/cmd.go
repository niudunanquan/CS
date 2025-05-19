package acme

import (
	"om/pkg/util"
	"os/exec"
	"runtime"
)

type LegoCommand struct {
	Binary string
}

func NewLegoCommand() LegoCommand {
	binary := util.RootDir + "lego"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	command := LegoCommand{
		Binary: binary,
	}

	return command
}

func (lc LegoCommand) ExecCommand(args ...string) *exec.Cmd {
	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, "--accept-tos", "--path", util.RootDir+"acme")
	cmdArgs = append(cmdArgs, args...)

	return exec.Command(lc.Binary, cmdArgs...)
}
