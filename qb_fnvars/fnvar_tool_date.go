package qb_fnvars

import (
	"fmt"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

const ToolDate = "date"

// FnVarToolDate date|yyyy-MM-dd|upper
type FnVarToolDate struct{}

func (instance *FnVarToolDate) Name() string {
	return ToolDate
}

func (instance *FnVarToolDate) Solve(token string, context ...interface{}) (interface{}, error) {
	tags := SplitToken(token)
	if len(tags) > 1 {
		return instance.solve(tags[1:])
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// date|yyyy-MM-dd|upper
// date|timestamp|add|day|30
func (instance *FnVarToolDate) solve(args []string) (interface{}, error) {
	pattern := qb_utils.Arrays.GetAt(args, 0, "yyyyMMdd").(string)
	lpattern := strings.ToLower(pattern)
	now := time.Now()
	var output string
	switch lpattern {
	case "iso", "rfc3339":
		output = now.Format(time.RFC3339)
	case "unix":
		output = now.Format(time.UnixDate)
	case "ruby":
		output = now.Format(time.RubyDate)
	case "timestamp", "timestamp_s":
		output = fmt.Sprintf("%d", now.Unix())
	case "timestamp_m":
		output = fmt.Sprintf("%d", now.Unix()*1000)
	default:
		output = qb_utils.Formatter.FormatDate(now, pattern)
	}

	transform := qb_utils.Arrays.GetAt(args, 1, "").(string)
	if len(transform) > 0 {
		switch transform {
		case "add", "sub":
			output = calculate(transform, args, output)
		default:
			output = Transform(output, transform)
		}
	}

	return output, nil
}

func calculate(op string, args []string, output string) string {
	timestamp := qb_utils.Convert.ToInt64Def(output, 0) // seconds
	val := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(args, 2, 0))
	um := qb_utils.Arrays.GetAt(args, 3, "second").(string)
	if timestamp > 0 && len(um) > 0 && val > 0 {
		if op == "sub" {
			val = val * -1
		}
		switch um {
		case "millisecond":
			timestamp = timestamp + int64(val/1000)
		case "second":
			timestamp = timestamp + int64(val)
		case "minute":
			timestamp = timestamp + int64(60*val)
		case "hour":
			timestamp = timestamp + int64(60*60*val)
		case "day":
			timestamp = qb_utils.Dates.AddDays(timestamp, val).Unix()
		case "week":
			timestamp = qb_utils.Dates.AddWeeks(timestamp, val).Unix()
		case "month":
			timestamp = qb_utils.Dates.AddMonths(timestamp, val).Unix()
		case "year":
			timestamp = qb_utils.Dates.AddYears(timestamp, val).Unix()
		default:
			// seconds
		}
		output = fmt.Sprintf("%d", timestamp)
	}
	return output
}
