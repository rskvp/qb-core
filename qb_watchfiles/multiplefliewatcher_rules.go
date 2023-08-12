package qb_watchfiles

import (
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

const (
	RulePDF  = "pdf"
	RuleTEXT = "text"
	RuleIMG  = "img"
)

type MultipleFileWatcherRules struct {
	Rules []*MultipleFileWatcherRule `json:"rules"`
}

func NewMultipleFileWatcherRules() (instance *MultipleFileWatcherRules) {
	instance = new(MultipleFileWatcherRules)
	instance.Rules = make([]*MultipleFileWatcherRule, 0)

	// add default
	instance.Rules = append(instance.Rules, &MultipleFileWatcherRule{
		FilePattern:   "*.pdf",
		FileValidator: RulePDF,
	})
	instance.Rules = append(instance.Rules, &MultipleFileWatcherRule{
		FilePattern:   "*.txt",
		FileValidator: RuleTEXT,
	})
	instance.Rules = append(instance.Rules, &MultipleFileWatcherRule{
		FilePattern:   "*.jpg, *.gif, *.png, *.bmp",
		FileValidator: RuleIMG,
	})

	return
}

func (instance *MultipleFileWatcherRules) String() string {
	return qb_utils.JSON.Stringify(instance)
}

func (instance *MultipleFileWatcherRules) IsValid(filename string) bool {
	if nil != instance {
		matched, validated, e := instance.Match(filename)
		notValid := nil != e || (matched && !validated)
		return !notValid
	}
	return false
}

func (instance *MultipleFileWatcherRules) Match(filename string) (match, validated bool, err error) {
	if nil != instance && len(instance.Rules) > 0 {
		for _, rule := range instance.Rules {
			if len(rule.FilePattern) > 0 && len(rule.FileValidator) > 0 {
				match, validated, err = rule.Match(filename)
				if match {
					break
				}
			}
		}
	}
	return
}

type MultipleFileWatcherRule struct {
	FilePattern   string `json:"file-pattern"`
	FileValidator string `json:"file-validator"` // "func", "text", "pdf", content regexp
}

func (instance *MultipleFileWatcherRule) String() string {
	return qb_utils.JSON.Stringify(instance)
}

func (instance *MultipleFileWatcherRule) Match(filename string) (match, validated bool, err error) {
	if nil != instance {
		if len(filename) > 0 {
			patterns := qb_utils.Strings.SplitTrimSpace(instance.FilePattern, ",")
			for _, pattern := range patterns {
				match = qb_utils.Paths.PatternMatchBase(filename, pattern)
				if match {
					validated, err = validate(instance.FileValidator, filename)
					return
				}
			}
		} else {
			err = qb_utils.Errors.Prefix(ErrorSystemPanic, "Missing file name: ")
		}
	} else {
		err = ErrorSystemPanic
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	v a l i d a t e
// ---------------------------------------------------------------------------------------------------------------------

func validate(validator string, filename string) (success bool, err error) {
	switch validator {
	case RulePDF:
		success, err = validatePDF(filename)
	case RuleTEXT:
		success, err = validateTEXT(filename)
	case RuleIMG:
		success, err = validateIMG(filename)
	default:
		err = qb_utils.Errors.Prefix(ErrorSystemPanic, fmt.Sprintf("Validator '%s' not supported: ", validator))
	}
	return
}

func validatePDF(filename string) (success bool, err error) {
	b, e := qb_utils.IO.ReadBytesFromFile(filename)
	if nil != e {
		err = e
		return
	}
	if len(b) > 6 {
		EOF := string(b[len(b)-5:])
		success = strings.Index(EOF, "EOF") > -1
	}
	return
}

func validateTEXT(filename string) (success bool, err error) {
	b, e := qb_utils.IO.ReadBytesFromFile(filename)
	if nil != e {
		err = e
		return
	}
	success = len(b) > 0
	return
}

func validateIMG(filename string) (success bool, err error) {
	b, e := qb_utils.IO.ReadBytesFromFile(filename)
	if nil != e {
		err = e
		return
	}
	success = len(b) > 3 // more than 3 bytes
	return
}
