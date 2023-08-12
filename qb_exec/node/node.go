package node

import (
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

var (
	Node *NodeHelper

	nodeCommand = "node"
)

const wpName = "node"
const fsName = "./.node"

type NodeHelper struct {
}

func init() {
	Node = new(NodeHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *NodeHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *NodeHelper) Version() (version string, err error) {
	program := instance.NewExec()
	version, err = program.Version()
	if nil == err {
		// version 2.23.0
		tokens := strings.Split(version, " ")
		if len(tokens) == 3 {
			version = tokens[2]
		}
	}
	return
}

// NewExec
// Creates new Solana command with default password
func (instance *NodeHelper) NewExec(args ...string) *NodeExec {
	if nil != instance {
		filename := "" // js file to execute
		if len(args) == 1 {
			filename = qb_utils.Convert.ToString(args[0])
		}
		return NewNodeExec(nodeCommand, nil, filename)
	}
	return nil
}
