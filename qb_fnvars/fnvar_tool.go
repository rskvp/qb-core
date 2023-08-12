package qb_fnvars

import (
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

const (
	prefix    = "<var>"
	suffix    = "</var>"
	separator = "|"
)

type FnVarTool interface {
	Name() string
	Solve(token string, context ...interface{}) (interface{}, error)
}

func GetTokens(text string) []string {
	return qb_utils.Regex.TagsBetweenTrimStrings(text+" ", prefix, suffix)
}

func SplitToken(token string) []string {
	clean := strings.Replace(strings.Replace(token, suffix, "", 1), prefix, "", 1)
	return strings.Split(clean, separator)
}

// ParseTools  a text and extract a map of all tools supported
func ParseTools(text string, registered map[string]FnVarTool) map[string]FnVarTool {
	m := make(map[string]FnVarTool)
	tokens := GetTokens(text)
	for _, token := range tokens {
		tool := GetTool(token, registered)
		if nil != tool {
			m[token] = tool
		}
	}
	return m
}

// GetTool retrieve a registered tool parsing token (i.e. $rnd:alphanumeric:6)
func GetTool(token string, registered map[string]FnVarTool) FnVarTool {
	tags := SplitToken(token)
	if len(tags) > 0 {
		name := tags[0]
		if v, b := registered[name]; b {
			return v
		}
	}
	return nil
}

func Transform(text string, transform string) string {
	switch transform {
	case "upper":
		return strings.ToUpper(text)
	case "lower":
		return strings.ToLower(text)
	default:
		return text
	}
}
