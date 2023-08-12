package nvm

import (
	_ "embed"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

var (
	Nvm *NvmHelper

	//go:embed tpl_nvm.sh
	tplNvm string
	//go:embed tpl_nvms.sh
	tplNvms string
	//go:embed tpl_node.sh
	tplNode string
	//go:embed tpl_npm.sh
	tplNpm string
	//go:embed tpl_npx.sh
	tplNpx string

	root = "./.nvm"
)

type NvmHelper struct {
}

func init() {
	Nvm = new(NvmHelper)

	userDir, err := qb_utils.Paths.UserHomeDir()
	if nil == err {
		root = qb_utils.Paths.Concat(userDir, ".nvm")
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *NvmHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *NvmHelper) Version() (version string, err error) {
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
// Creates new nvm command
func (instance *NvmHelper) NewExec() *NvmExec {
	if nil != instance {
		return NewExec(nil)
	}
	return nil
}
