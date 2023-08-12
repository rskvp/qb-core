package qb_exp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

//bracketType sorts characters into either opening, closing, or not a bracket
type bracketType int

const (
	openBracket bracketType = iota
	closeBracket
	notABracket
)

//bracketPairs are the matching pairs of brackets
var pairs = map[rune]rune{'{': '}', '[': ']', '(': ')'}

var operators = []string{"|", "&", "+", "-", "*", "%"}

//----------------------------------------------------------------------------------------------------------------------
//	Statement
//----------------------------------------------------------------------------------------------------------------------

type Statement struct {
	Uid      string
	Text     string
	Operator string
	Index    int
	Count    int
	Inner    int          // inner level
	Next     []*Statement // on the right of statement
}

func (instance *Statement) String() string {
	if instance.IsLeaf() {
		if instance.IsLast() {
			return instance.Text
		}
		return fmt.Sprintf("%s%s", instance.Text, instance.Operator)
	} else {
		group := instance.IsBlock()
		var s bytes.Buffer
		if group {
			s.WriteString("(")
		}
		for _, statement := range instance.Next {
			s.WriteString(statement.String())
		}
		if group {
			s.WriteString(")")
		}
		s.WriteString(instance.Next[len(instance.Next)-1].Operator)

		return s.String()
	}
}

// IsBlock if true, this expression is self-contained in brackets
func (instance *Statement) IsBlock() bool {
	return instance.Inner == 0 || instance.IsFirst() || instance.IsLast()
}

func (instance *Statement) IsLeaf() bool {
	return len(instance.Next) == 0 // no children
}

// IsBranch return true if the statement has at least one leaf
func (instance *Statement) IsBranch() bool {
	if len(instance.Next) > 0 {
		for _, statement := range instance.Next {
			if statement.IsLeaf() {
				return true
			}
		}
	} else {
		return instance.IsFirst() && instance.Count == 1
	}
	return false
}

func (instance *Statement) IsFirst() bool {
	return instance.Index == 0
}

func (instance *Statement) IsLast() bool {
	return instance.Index == instance.Count-1
}

func (instance *Statement) HasChildren() bool {
	return len(instance.Next) > 0 // some children
}

func (instance *Statement) Explain() string {
	var sb bytes.Buffer
	explain(0, instance, &sb)
	return sb.String()
}

func (instance *Statement) WalkThrough(callback func(level, index int, statement *Statement)) {
	walkThrough(0, 0, instance, callback)
}

//----------------------------------------------------------------------------------------------------------------------
//	StatementGroup
//----------------------------------------------------------------------------------------------------------------------

type StatementGroup []*Statement

func (instance *StatementGroup) String() string {
	var s bytes.Buffer
	for _, statement := range *instance {
		s.WriteString(statement.String())
	}
	return s.String()
}

func (instance *StatementGroup) Parse(text string) error {
	if len(text) > 0 {
		if !IsBalanced(text) {
			return errors.New(fmt.Sprintf("Text is not balanced. Please, check parentheses are properly coupled: '%s'", text))
		}

		// add statement
		*instance = append(*instance, parseSibling("", 0, text)...)
		return nil
	}
	return errors.New("invalid text: text should not be empty")
}

func (instance *StatementGroup) Explain() string {
	var s bytes.Buffer
	for _, statement := range *instance {
		explain(0, statement, &s)
	}
	return s.String()
}

func (instance *StatementGroup) WalkThrough(callback func(level, index int, statement *Statement)) {
	for i, statement := range *instance {
		walkThrough(0, i, statement, callback)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func SetOperators(values []string) {
	operators = values
}

func GetOperators() []string {
	return operators
}

/*
ParseStatement group statement into nested arrays
@text "variable + variable - anothervariable"
*/
func ParseStatement(text string) (*StatementGroup, error) {
	response := &StatementGroup{}
	err := response.Parse(text)
	if nil != err {
		return nil, err
	}
	return response, nil
}

func IsStatement(phrase string) bool {
	for _, c := range phrase {
		if b, _ := isOperator(string(c)); b {
			return true
		}
		bt := getBracketType(c)
		if bt != notABracket {
			return true
		}
	}

	return false
}

/*IsBalanced determines if a strings has balanced brackets*/
func IsBalanced(phrase string) bool {
	var queue []rune
	for _, v := range phrase {
		switch getBracketType(v) {
		case openBracket:
			queue = append(queue, pairs[v])
		case closeBracket:
			if 0 < len(queue) && queue[len(queue)-1] == v {
				queue = queue[:len(queue)-1]
			} else {
				return false
			}
		}
	}
	return len(queue) == 0
}

// SplitSibling get "A + (....) + B + (...)"
// and return ["A +", "(...)", "+ B", "(...)"]
func SplitSibling(phrase string) []string {
	phrase = prepareText(phrase, operators[0])
	response := make([]string, 0)
	queue := 0
	closeBuffer := false
	var buf bytes.Buffer
	for _, v := range phrase {
		switch getBracketType(v) {
		case openBracket:
			if queue == 0 {
				if buf.Len() > 0 {
					response = append(response, buf.String())
					buf.Reset()
				}
			} else {
				buf.WriteRune(v)
			}
			queue++
		case closeBracket:
			queue--
			if queue == 0 {
				// response = append(response, buf.String())
				// buf.Reset()
				closeBuffer = true
			} else {
				buf.WriteRune(v)
			}
		default:
			buf.WriteRune(v)
			if closeBuffer {
				closeBuffer = false
				if buf.Len() > 0 {
					response = append(response, buf.String())
					buf.Reset()
				}
			}
		}
	}
	if buf.Len() > 0 {
		response = append(response, buf.String())
		buf.Reset()
	}
	return response
}

func Normalize(phrase string) string {
	return prepareText(phrase, operators[0])
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

/*getBracketType determines the type of bracket*/
func getBracketType(char rune) bracketType {
	for k, v := range pairs {
		switch char {
		case k:
			return openBracket
		case v:
			return closeBracket
		}
	}
	return notABracket
}

func isOperator(char string) (bool, string) {
	for _, v := range operators {
		if v == char {
			return true, v
		}
	}
	return false, ""
}

func hasBrackets(phrase string) bool {
	for k, v := range pairs {
		if strings.Index(phrase, string(k)) > -1 || strings.Index(phrase, string(v)) > -1 {
			return true
		}
	}
	return false
}

func prepareText(text, defOperator string) string {
	// prepare text
	text = qb_utils.Strings.RemoveDuplicateSpaces(text)
	// adjust operators removing spaces
	for _, v := range operators {
		text = strings.ReplaceAll(text, " "+v+" ", v)
		text = strings.ReplaceAll(text, " "+v, v)
		text = strings.ReplaceAll(text, v+" ", v)
	}
	// add default operator replacing space
	text = strings.ReplaceAll(text, " ", defOperator)

	return text
}

func uid(rootUid string, inner, index int) string {
	var sb bytes.Buffer
	if len(rootUid) > 0 {
		sb.WriteString(fmt.Sprintf("%s.", rootUid))
	}
	sb.WriteString(fmt.Sprintf("%v", index+1))
	return sb.String()
}

func parseSibling(rootUid string, inner int, phrase string) []*Statement {
	response := make([]*Statement, 0)
	siblings := SplitSibling(phrase)
	if len(siblings) > 0 {
		for i, text := range siblings {
			if hasBrackets(text) {
				statement := new(Statement)
				statement.Uid = uid(rootUid, inner, i)
				statement.Index = i
				statement.Count = len(siblings)
				statement.Inner = inner
				statement.Text = text
				statement.Next = make([]*Statement, 0)
				statement.Next = parseSibling(statement.Uid, inner+1, text)
				// add statement
				response = append(response, statement)
			} else {
				statements, seps := qb_utils.Strings.SplitAfter(text, strings.Join(operators, ""))
				if len(statements) == 1 {
					statement := new(Statement)
					statement.Uid = uid(rootUid, inner, i)
					statement.Index = i
					statement.Count = len(siblings)
					statement.Inner = inner
					statement.Text = statements[0]
					statement.Operator = qb_utils.Arrays.GetAt(seps, 0, "").(string)
					statement.Next = make([]*Statement, 0)
					// add statement
					response = append(response, statement)
				} else {
					master := new(Statement)
					master.Uid = uid(rootUid, inner, i)
					master.Text = text
					master.Index = i
					master.Count = len(siblings)
					master.Inner = inner
					master.Operator = ""
					master.Next = make([]*Statement, 0)
					response = append(response, master) // add statement
					for ii, st := range statements {
						statement := new(Statement)
						statement.Uid = uid(master.Uid, inner+1, ii)
						statement.Index = ii
						statement.Count = len(statements)
						statement.Inner = inner + 1
						statement.Text = st
						statement.Operator = qb_utils.Arrays.GetAt(seps, ii, "").(string)
						statement.Next = make([]*Statement, 0)
						// add to master
						master.Next = append(master.Next, statement)
					}
				}
			}
		}
	} else {
		// single statement
		statement := new(Statement)
		statement.Uid = uid(rootUid, inner, 0)
		statement.Text = phrase
		statement.Operator = ""
		statement.Next = make([]*Statement, 0)
		response = append(response, statement)
	}
	return response
}

func explain(level int, statement *Statement, out *bytes.Buffer) {
	indent := qb_utils.Strings.Repeat("\t", level)
	if out.Len() > 0 {
		out.WriteString("\n")
	}
	out.WriteString(fmt.Sprintf("%s%s) %s [is-block:%v, is-branch:%v]", indent, statement.Uid, statement.String(), statement.IsBlock(), statement.IsBranch()))
	if statement.HasChildren() {
		for _, child := range statement.Next {
			explain(level+1, child, out)
		}
	}
}

func walkThrough(level, index int, statement *Statement, callback func(level, index int, statement *Statement)) {
	if nil == callback {
		panic("missing callback parameter.")
	}
	// call root
	callback(level, index, statement)
	if statement.HasChildren() {
		for i, child := range statement.Next {
			// recursive
			walkThrough(level+1, i, child, callback)
		}
	}
}
