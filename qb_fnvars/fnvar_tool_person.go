package qb_fnvars

import (
	"github.com/rskvp/qb-core/qb_generators/genusers"
	"github.com/rskvp/qb-core/qb_utils"
)

const ToolUser = "user"

// FnVarToolUser person|name,mobile,mail|3
type FnVarToolUser struct {
	dataRoot string
}

func (instance *FnVarToolUser) Name() string {
	return ToolUser
}

func (instance *FnVarToolUser) Solve(token string, context ...interface{}) (interface{}, error) {
	tags := SplitToken(token)
	if len(tags) > 0 {
		return instance.solve(tags[1:])
	}
	return nil, nil
}

func (instance *FnVarToolUser) SetDataRoot(value string) {
	instance.dataRoot = value
}

func (instance *FnVarToolUser) GetDataRoot() string {
	return instance.dataRoot
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FnVarToolUser) getGenUserEngine() *genusers.GenUsersEngine {
	return genusers.GenUsers.NewEngine(instance.dataRoot)
}

func (instance *FnVarToolUser) solve(args []string) (interface{}, error) {
	response := make([]map[string]interface{}, 0)
	options := qb_utils.Arrays.GetAt(args, 0, "").(string)
	count := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(args, 1, "1").(string))
	fields := qb_utils.Strings.SplitTrimSpace(options, ",")

	generator := instance.getGenUserEngine()
	users, err := generator.Generate(count)
	if nil != err {
		return nil, err
	}
	for _, user := range users {
		response = append(response, user.Map(fields...))
	}

	if len(response) == 1 {
		item := response[0]
		if len(fields) > 1 || len(fields) == 0 || fields[0] == "*" || fields[0] == "" {
			return item, nil
		}
		return item[fields[0]], nil
	}
	return response, nil
}
