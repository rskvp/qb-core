package qb_utils

import (
	"bytes"
	"errors"
	"fmt"
	templateHtml "html/template"
	"math"
	"strconv"
	"strings"
	templateText "text/template"
	"time"
)

type FormatterHelper struct {
}

var Formatter *FormatterHelper

func init() {
	Formatter = new(FormatterHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	B Y T E S
//----------------------------------------------------------------------------------------------------------------------

func (instance *FormatterHelper) FormatBytes(i interface{}) string {
	n := uint64(Convert.ToInt64(i))
	return instance.fmtBytes(n)
}

func (instance *FormatterHelper) fmtBytes(b uint64) string {
	const unit = 1024
	if b < Kb {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

//----------------------------------------------------------------------------------------------------------------------
//	T E M P L A T E
//----------------------------------------------------------------------------------------------------------------------

func (instance *FormatterHelper) Render(s string, data interface{}, allowEmpty bool) string {
	response := s
	fields := Regex.MatchBetween(s, 0, "{{", "}}", "")
	for _, field := range fields {
		fieldName := strings.Trim(field, " \n.")
		value := Reflect.Get(data, fieldName)
		if t, ok := value.(string); ok {
			if allowEmpty || len(t) > 0 {
				response = strings.ReplaceAll(response, fmt.Sprintf("{{%s}}", field), t)
			}
		}
	}
	return response
}

func (instance *FormatterHelper) Merge(textOrHtml string, data interface{}) (string, error) {
	if Regex.IsHTML(textOrHtml) {
		// safe HTML merge (escape HTML tags)
		return mergeHtml(textOrHtml, data, false)
	} else {
		return mergeText(textOrHtml, data, false)
	}
}

func (instance *FormatterHelper) MergeKeepFields(textOrHtml string, data interface{}) (string, error) {
	if Regex.IsHTML(textOrHtml) {
		// safe HTML merge (escape HTML tags)
		return mergeHtml(textOrHtml, data, true)
	} else {
		return mergeText(textOrHtml, data, true)
	}
}

func (instance *FormatterHelper) MergeText(text string, data interface{}) (string, error) {
	return mergeText(text, data, false)
}

func (instance *FormatterHelper) MergeTextKeepFields(text string, data interface{}) (string, error) {
	return mergeText(text, data, true)
}

func (instance *FormatterHelper) MergeHtml(html string, data interface{}) (string, error) {
	return mergeHtml(html, data, false)
}

func (instance *FormatterHelper) MergeHtmlKeepFields(html string, data interface{}) (string, error) {
	return mergeHtml(html, data, true)
}

func (instance *FormatterHelper) MergeDef(textOrHtml string, data interface{}, def string) (string, error) {
	merged, err := instance.Merge(textOrHtml, data)
	if len(merged) > 0 {
		return merged, err
	}
	return def, err
}

func (instance *FormatterHelper) MergeDefKeepFields(textOrHtml string, data interface{}, def string) (string, error) {
	merged, err := instance.MergeKeepFields(textOrHtml, data)
	if len(merged) > 0 {
		return merged, err
	}
	return def, err
}

func (instance *FormatterHelper) MergeTextDef(text string, data interface{}, def string) (string, error) {
	merged, err := mergeText(text, data, false)
	if len(merged) > 0 {
		return merged, err
	}
	return def, err
}

func (instance *FormatterHelper) MergeHtmlDef(text string, data interface{}, def string) (string, error) {
	merged, err := mergeHtml(text, data, false)
	if len(merged) > 0 {
		return merged, err
	}
	return def, err
}

func (instance *FormatterHelper) MergeRecover(textOrHtml string, data interface{}, def string, callback func(string, error)) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			e := fmt.Sprintf("%s", r)
			if nil != callback {
				callback(def, errors.New(e))
			}
		}
	}()
	response, err := instance.MergeDef(textOrHtml, data, def)
	if nil != callback {
		callback(response, err)
	}
}

func (instance *FormatterHelper) MergeTextRecover(text string, data interface{}, def string, callback func(string, error)) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			e := fmt.Sprintf("%s", r)
			if nil != callback {
				callback(def, errors.New(e))
			}
		}
	}()
	response, err := instance.MergeTextDef(text, data, def)
	if nil != callback {
		callback(response, err)
	}
}

func (instance *FormatterHelper) MergeHtmlRecover(text string, data interface{}, def string, callback func(string, error)) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			e := fmt.Sprintf("%s", r)
			if nil != callback {
				callback(def, errors.New(e))
			}
		}
	}()
	response, err := instance.MergeHtmlDef(text, data, def)
	if nil != callback {
		callback(response, err)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	D A T E
//----------------------------------------------------------------------------------------------------------------------

func (instance *FormatterHelper) FormatDate(dt time.Time, pattern string) string {
	return Dates.FormatDate(dt, pattern)
}

func (instance *FormatterHelper) ParseDate(dt string, pattern string) (time.Time, error) {
	return Dates.ParseDate(dt, pattern)
}

//----------------------------------------------------------------------------------------------------------------------
//	M A P
//----------------------------------------------------------------------------------------------------------------------

func (instance *FormatterHelper) FormatMap(i interface{}) string {
	m := Convert.ToMap(i)
	if nil != m {
		return printMap(m, 0)
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	N U M B E R S
//----------------------------------------------------------------------------------------------------------------------

var renderFloatPrecisionMultipliers = [10]float64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
}

var renderFloatPrecisionRounders = [10]float64{
	0.5,
	0.05,
	0.005,
	0.0005,
	0.00005,
	0.000005,
	0.0000005,
	0.00000005,
	0.000000005,
	0.0000000005,
}

func (instance *FormatterHelper) FormatInteger(i interface{}, pattern string) string {
	if len(pattern) == 0 {
		pattern = "#,###."
	}
	n := Convert.ToInt64(i)
	return renderFloat(pattern, float64(n))
}

func (instance *FormatterHelper) FormatFloat(i interface{}, pattern string) string {
	if len(pattern) == 0 {
		pattern = "#,###.##"
	}
	n := Convert.ToFloat64(i)
	return renderFloat(pattern, n)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func mergeText(text string, data interface{}, keepFields bool) (string, error) {
	if nil == data {
		data = struct{}{}
	}
	buff := bytes.NewBufferString("")
	if keepFields {
		text = keepTemplateFields(text)
	}
	t, err := templateText.New("template").Parse(strings.Trim(text, "\n"))
	if nil != err {
		return "", err
	}
	err = t.Execute(buff, data)
	if nil != err {
		return "", err
	}
	return buff.String(), err
}

func mergeHtml(html string, data interface{}, keepFields bool) (string, error) {
	if nil == data {
		data = struct{}{}
	}
	buff := bytes.NewBufferString("")
	if keepFields {
		html = keepTemplateFields(html)
	}
	t, err := templateHtml.New("template").Parse(strings.Trim(html, "\n"))
	if nil != err {
		return "", err
	}
	err = t.Execute(buff, data)
	if nil != err {
		return "", err
	}
	return buff.String(), err
}

func keepTemplateFields(text string) (response string) {
	response = text
	fields := Regex.MatchBetween(response, 0, "{{", "}}", "")
	for _, field := range fields {
		value := fmt.Sprintf("{{if%s}}{{%s}}{{else}}{{\"{{%s}}\"}}{{end}}", field, field, field)
		response = strings.ReplaceAll(response, fmt.Sprintf("{{%s}}", field), value)
	}
	return
}

func printMap(m map[string]interface{}, level int) string {
	var buf bytes.Buffer
	prefix := ""
	if level > 0 {
		prefix = Strings.FillLeft("", level, ' ')
	}
	prefix = strings.ReplaceAll(prefix, " ", "\t")
	for k, v := range m {
		if mm, b := v.(map[string]interface{}); b {
			buf.WriteString(fmt.Sprintf("%v %v:\n", prefix, k))
			buf.WriteString(printMap(mm, level+1))
		} else {
			buf.WriteString(fmt.Sprintf("%v %v: %v\n", prefix, k, v))
		}
	}
	return buf.String()
}

/*
*

	Examples of format strings, given n = 12345.6789:
	"#,###.##" => "12,345.67"
	"#,###." => "12,345"
	"#,###" => "12345,678"
	"#\u202F###,##" => "12 345,67"
	"#.###,###### => 12.345,678900
	"" (aka default format) => 12,345.67
	The highest precision allowed is 9 digits after the decimal symbol.

*
*/
func renderFloat(format string, n float64) string {
	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	if math.IsNaN(n) {
		return "NaN"
	}
	if n > math.MaxFloat64 {
		return "Infinity"
	}
	if n < -math.MaxFloat64 {
		return "-Infinity"
	}

	// default format
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = ""

		// collect indices of meaningful formatting directives
		formatDirectiveChars := []rune(format)
		formatDirectiveIndices := make([]int, 0)
		for i, char := range formatDirectiveChars {
			if char != '#' && char != '0' {
				formatDirectiveIndices = append(formatDirectiveIndices, i)
			}
		}

		if len(formatDirectiveIndices) > 0 {
			// Directive at index 0:
			//   Must be a '+'
			//   Raise an error if not the case
			// index: 0123456789
			//        +0.000,000
			//        +000,000.0
			//        +0000.00
			//        +0000
			if formatDirectiveIndices[0] == 0 {
				if formatDirectiveChars[formatDirectiveIndices[0]] != '+' {
					panic("renderFloat(): invalid positive sign directive")
				}
				positiveStr = "+"
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// Two directives:
			//   First is thousands separator
			//   Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatDirectiveIndices) == 2 {
				if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
					panic("renderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				}
				thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// One directive:
			//   Directive is decimal separator
			//   The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatDirectiveIndices) == 1 {
				decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
			}
		}
	}

	// generate sign part
	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.Itoa(int(intf))

	// add thousand separator if required
	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if precision == 0 {
		return signStr + intStr
	}

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr
}
