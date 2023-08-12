package yarn

import (
	"strings"
)

var (
	Yarn *YarnHelper

	yarnCommand = "yarn"
)

const wpName = "yarn"
const fsName = "./.yarn"

type YarnHelper struct {
}

func init() {
	Yarn = new(YarnHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *YarnHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *YarnHelper) Version() (version string, err error) {
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
// Creates new yarn exec
func (instance *YarnHelper) NewExec() *YarnExec {
	if nil != instance {
		return NewExec(yarnCommand, nil)
	}
	return nil
}
