package git

import (
	"strings"
)

var (
	Git *GitHelper

	gitCommand = "git"
)

const wpName = "git"
const fsName = "./.git"

type GitHelper struct {
}

func init() {
	Git = new(GitHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GitHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *GitHelper) Version() (version string, err error) {
	program := instance.NewExec()
	version, err = program.Version()
	if nil == err {
		// git version 2.23.0
		tokens := strings.Split(version, " ")
		if len(tokens) == 3 {
			version = tokens[2]
		}
	}
	return
}

// NewExec
// Creates new Solana command with default password
func (instance *GitHelper) NewExec() *GitExec {
	if nil != instance {
		return NewExec(gitCommand, nil)
	}
	return nil
}
