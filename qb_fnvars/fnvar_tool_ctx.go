package qb_fnvars

import "github.com/rskvp/qb-core/qb_utils"

const (
	ToolCtx = "ctx"
)

// FnVarToolCtx ctx|var1|4, ctx|var2
type FnVarToolCtx struct{}

func (instance *FnVarToolCtx) Name() string {
	return ToolCtx
}

func (instance *FnVarToolCtx) Solve(token string, context ...interface{}) (interface{}, error) {
	tags := SplitToken(token)
	if len(tags) > 1 {
		return instance.solve(tags[1:], context...)
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FnVarToolCtx) solve(args []string, context ...interface{}) (interface{}, error) {
	if len(context) == 1 && len(args) > 0 {
		ctx := context[0]
		return instance.get(ctx, args), nil
	}
	return "", nil
}

func (instance *FnVarToolCtx) get(obj interface{}, tokens []string) interface{} {
	response := obj
	for _, n := range tokens {
		response = qb_utils.Reflect.Get(response, n)
	}
	return response
}
