package npm

import (
	"strings"
)

var (
	Npm *NpmHelper

	npmCommand = "npm"
)

const wpName = "npm"
const fsName = "./.npm"

type NpmHelper struct {
}

func init() {
	Npm = new(NpmHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *NpmHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *NpmHelper) Version() (version string, err error) {
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
// Creates new npm command
func (instance *NpmHelper) NewExec() *NpmExec {
	if nil != instance {
		return NewExec(npmCommand, nil)
	}
	return nil
}
