// +build linux

package qb_exec

import "os/exec"

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	OPEN_FILE_COMMAND = "xdg-open"
)

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func openFileCommand(args ...string) *exec.Cmd {
	return exec.Command(OPEN_FILE_COMMAND, args...)
}
