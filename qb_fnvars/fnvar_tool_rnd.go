package qb_fnvars

import (
	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_utils"
)

const (
	ToolRnd = "rnd"
)

// FnVarToolRnd rnd|chars|4|upper
type FnVarToolRnd struct{}

func (instance *FnVarToolRnd) Name() string {
	return ToolRnd
}

func (instance *FnVarToolRnd) Solve(token string, context ...interface{}) (interface{}, error) {
	tags := SplitToken(token)
	if len(tags) > 1 {
		return instance.solve(tags[1:])
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FnVarToolRnd) solve(args []string) (interface{}, error) {
	mode := qb_utils.Arrays.GetAt(args, 0, "numeric")
	var output string
	switch mode {
	case "numeric":
		length := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(args, 1, "6"))
		output = qb_rnd.Rnd.RndDigits(length)
	case "range", "between":
		param := qb_utils.Convert.ToString(qb_utils.Arrays.GetAt(args, 1, "0-10"))
		tokens := qb_utils.Strings.Split(param, "-.,:;")
		from := qb_utils.Convert.ToInt64(tokens[0])
		to := from + 5
		if len(tokens) > 1 {
			to = qb_utils.Convert.ToInt64(tokens[1])
		}
		output = qb_utils.Convert.ToString(qb_rnd.Rnd.Between(from, to))
	case "alphanumeric", "chars":
		length := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(args, 1, "6"))
		transform := qb_utils.Arrays.GetAt(args, 2, "")
		switch transform {
		case "upper":
			output = qb_rnd.Rnd.RndCharsUpper(length)
		case "lower":
			output = qb_rnd.Rnd.RndCharsLower(length)
		default:
			output = qb_rnd.Rnd.RndChars(length)
		}
	case "guid":
		transform := qb_utils.Arrays.GetAt(args, 2, "").(string)
		output = qb_rnd.Rnd.Uuid()
		output = Transform(output, transform)
	case "id":
		transform := qb_utils.Arrays.GetAt(args, 2, "").(string)
		output = qb_rnd.Rnd.RndId()
		output = Transform(output, transform)
	default:
		length := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(args, 1, "6"))
		output = qb_rnd.Rnd.RndDigits(length)
	}
	return output, nil
}
