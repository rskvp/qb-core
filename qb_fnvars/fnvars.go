package qb_fnvars

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

type FnVarsHelper struct {
}

var FnVars *FnVarsHelper

func init() {
	FnVars = new(FnVarsHelper)
}

func (instance *FnVarsHelper) NewEngine() *FnVarsEngine {
	response := new(FnVarsEngine)
	response.registered = make(map[string]FnVarTool)

	response.init()

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	FnVarsEngine
//----------------------------------------------------------------------------------------------------------------------

type FnVarsEngine struct {
	registered map[string]FnVarTool
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *FnVarsEngine) Register(tool FnVarTool) {
	if nil != tool {
		name := tool.Name()
		if _, b := instance.registered[name]; !b {
			instance.registered[name] = tool
		}
	}
}

func (instance *FnVarsEngine) GetByName(toolName string) FnVarTool {
	if nil != instance {
		if v, ok := instance.registered[toolName]; ok {
			return v
		}
	}
	return nil
}

func (instance *FnVarsEngine) Solve(input interface{}, context ...interface{}) (interface{}, error) {
	if t, b := input.(string); b {
		return instance.SolveText(t, context...)
	} else if m, b := input.(map[string]string); b {
		return instance.SolveMap(m, context...)
	} else if a, b := input.([]string); b {
		return instance.SolveArray(a)
	}
	return nil, errors.New(fmt.Sprintf("Not supported input: %v", input))
}

func (instance *FnVarsEngine) SolveText(text string, context ...interface{}) (string, error) {
	model := ParseTools(text, instance.registered)
	if len(model) > 0 {
		for k, v := range model {
			value, err := v.Solve(k, context...)
			if nil != err {
				return text, err
			}
			text = strings.ReplaceAll(text, k, qb_utils.Convert.ToString(value))
		}
	}
	return text, nil
}

func (instance *FnVarsEngine) SolveMap(m map[string]string, context ...interface{}) (map[string]string, error) {
	for k, v := range m {
		solved, err := instance.SolveText(v, context...)
		if nil != err {
			return m, err
		}
		m[k] = solved
	}
	return m, nil
}

func (instance *FnVarsEngine) SolveArray(a []string, context ...interface{}) ([]string, error) {
	for k, v := range a {
		solved, err := instance.SolveText(v, context...)
		if nil != err {
			return a, err
		}
		a[k] = solved
	}
	return a, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FnVarsEngine) init() {
	instance.Register(new(FnVarToolRnd))
	instance.Register(new(FnVarToolDate))
	instance.Register(new(FnVarToolUser))
	instance.Register(new(FnVarToolCtx))
}
