package node_ts

import (
	"strings"

	"github.com/rskvp/qb-core/qb_exec/npm"
	"github.com/rskvp/qb-core/qb_utils"
)

var (
	TsNode *TsNodeHelper

	tsnodeCommand = "ts-node"
)

const wpName = "node"
const fsName = "./.node"

type TsNodeHelper struct {
}

func init() {
	TsNode = new(TsNodeHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *TsNodeHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *TsNodeHelper) Version() (version string, err error) {
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

func (instance *TsNodeHelper) Install() (version string, err error) {
	version, err = instance.Version()
	if nil != err {
		// try to install
		npmExec := npm.Npm.NewExec()
		_, err = npmExec.Version()
		if nil != err {
			return
		}
		/*
			globally with TypeScript.
			npm install -g typescript
			npm install -g ts-node
			npm install -D tslib @types/node
		*/
		_, err = npmExec.ExecuteCommand("install", "-g", "typescript")
		if nil == err {
			_, err = npmExec.ExecuteCommand("install", "-g", "ts-node")
			if nil == err {
				_, err = npmExec.ExecuteCommand("install", "-D", "tslib", "@types/node")
				version, err = instance.Version()
			}
		}
	}
	return
}

// NewExec
// Creates new exec command with default password
func (instance *TsNodeHelper) NewExec(args ...string) *TsNodeExec {
	if nil != instance {
		filename := "" // js file to execute
		if len(args) == 1 {
			filename = qb_utils.Convert.ToString(args[0])
		}
		return NewTsNodeExec(tsnodeCommand, nil, filename)
	}
	return nil
}
