package qb_utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const escape = '\\'

type StringsHelper struct {
}

var Strings *StringsHelper

func init() {
	Strings = new(StringsHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *StringsHelper) RemoveDuplicateSpaces(text string) string {
	space := regexp.MustCompile(`\s+`) // [\s\p{Zs}]{2,}
	return space.ReplaceAllString(text, " ")
}

func (instance *StringsHelper) TrimSpaces(slice []string) {
	instance.Trim(slice, " ")
}

func (instance *StringsHelper) Trim(slice []string, trimVal string) {
	for i := range slice {
		slice[i] = strings.Trim(slice[i], trimVal)
	}
}

func (instance *StringsHelper) Clear(text string) string {
	var buf bytes.Buffer
	lines := strings.Split(text, "\n")
	count := 0
	for _, line := range lines {
		space := regexp.MustCompile(`\s+`)
		s := strings.TrimSpace(space.ReplaceAllString(line, " "))
		if len(s) > 0 {
			if count > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(strings.TrimSpace(s))
			count++
		}
	}
	return buf.String()
}

//Paginate Clear and paginate text on multiple rows.
// New rows only if prev char is a dot (.)
func (instance *StringsHelper) Paginate(text string) string {
	// split rows
	lines := strings.Split(instance.Clear(text), ".\n")
	var buf bytes.Buffer
	count := 0
	for _, line := range lines {
		if count > 0 {
			buf.WriteString(".\n")
		}
		buf.WriteString(strings.ReplaceAll(line, "\n", " "))
		count++
	}
	return buf.String()
}

func (instance *StringsHelper) Concat(params ...interface{}) string {
	result := ""
	for _, v := range params {
		result += Convert.ToString(v)
	}
	return result
}

func (instance *StringsHelper) ConcatSep(separator string, params ...interface{}) string {
	result := ""
	strParams := Convert.ToArrayOfString(params...)
	for _, value := range strParams {
		if len(result) > 0 {
			result += separator
		}
		result += value
	}
	return result
}

func (instance *StringsHelper) ConcatTrimSep(separator string, params ...interface{}) string {
	result := ""
	for _, v := range params {
		value := strings.TrimSpace(Convert.ToString(v))
		if len(value) > 0 {
			if len(result) > 0 {
				result += separator
			}
			result += value
		}
	}
	return result
}

func (instance *StringsHelper) Format(s string, params ...interface{}) string {
	return fmt.Sprintf(strings.Replace(s, "%s", "%v", -1), params...)
}

func (instance *StringsHelper) FormatValues(s string, params ...interface{}) string {
	return fmt.Sprintf(s, params...)
}

// Split using all rune in a string of separators
func (instance *StringsHelper) Split(s string, seps string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, sep := range seps {
			if r == sep {
				return true
			}
		}
		return false
	})
}

func (instance *StringsHelper) SplitAfter(s string, seps string) (tokens []string, separators []string) {
	tokens = strings.FieldsFunc(s, func(r rune) bool {
		for _, sep := range seps {
			if r == sep {
				separators = append(separators, string(sep))
				return true
			}
		}
		return false
	})
	return
}

func (instance *StringsHelper) SplitTrim(s string, seps string, cutset string) []string {
	data := instance.Split(s, seps)
	for i, item := range data {
		data[i] = strings.Trim(item, cutset)
	}
	return data
}

func (instance *StringsHelper) SplitTrimSpace(s string, seps string) []string {
	return instance.SplitTrim(s, seps, " ")
}

func (instance *StringsHelper) SplitAndGetAt(s string, seps string, index int) string {
	tokens := instance.Split(s, seps)
	if len(tokens) > index {
		return tokens[index]
	}
	return ""
}

func (instance *StringsHelper) SplitLast(s string, sep rune) []string {
	data := strings.Split(s, string(sep))
	if len(data) > 1 {
		return []string{strings.Join(data[:len(data)-1], string(sep)), data[len(data)-1]}
	}
	return data
}

// SplitQuoted splits a string, ignoring separators present inside quoted runs.  Separators
// cannot be escaped outside quoted runs, the escaping will be ignored.
//
// Quotes are preserved in result, but the separators are removed.
func (instance *StringsHelper) SplitQuoted(s string, sep rune, quote rune) []string {
	a := make([]string, 0, 8)
	quoted := false
	escaped := false
	p := 0
	for i, c := range s {
		if c == escape {
			// Escape can escape itself.
			escaped = !escaped
			continue
		}
		if c == quote {
			quoted = !quoted
			continue
		}
		escaped = false
		if !quoted && c == sep {
			a = append(a, s[p:i])
			p = i + 1
		}
	}

	if quoted && quote != 0 {
		// s contained an unterminated quoted-run, re-split without quoting.
		return instance.SplitQuoted(s, sep, rune(0))
	}

	return append(a, s[p:])
}

// Sub get a substring
// @param s string The string
// @param start int Start index
// @param end int End index
func (instance *StringsHelper) Sub(s string, start int, end int) string {
	runes := []rune(s) // convert in rune to handle all characters.
	if start < 0 || start > end {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}

	return string(runes[start:end])
}

func (instance *StringsHelper) SubBetween(s string, prefix, suffix string) string {
	start, end := instance.IndexStartEnd(s, prefix, suffix)
	if start > -1 && end > -1 {
		return instance.Sub(s, start, end)
	}
	return ""
}

func (instance *StringsHelper) Contains(s string, seps string) bool {
	for _, r := range seps {
		if strings.Index(s, string(r)) > -1 {
			return true
		}
	}
	return false
}

func (instance *StringsHelper) IndexStartEnd(s string, prefix, suffix string) (start, end int) {
	start = -1
	end = -1
	start = strings.Index(s, prefix)
	if start == -1 {
		return
	}

	lprefix := len(prefix)
	lsuffix := len(suffix)
	sub := s[start+lprefix:] // get substring starting from end of found text
	end = strings.Index(sub, suffix)
	if end > -1 {
		end += start + lprefix + lsuffix
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	n o r m a l i z a t i o n
//----------------------------------------------------------------------------------------------------------------------

func (instance *StringsHelper) Slugify(text string, replaces ...string) string {
	// remove duplicate spaces, carriage returns and tabs
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")
	if len(replaces) > 0 {
		for _, replace := range replaces {
			if len(replace) > 1 && strings.Index(replace, ":") > -1 {
				tokens := strings.Split(replace, ":")
				text = strings.ReplaceAll(text, tokens[0], tokens[1])
			} else {
				text = strings.ReplaceAll(text, " ", replace)
			}
		}
		return text
	} else {
		return strings.ReplaceAll(text, " ", "-")
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	Underscore
//----------------------------------------------------------------------------------------------------------------------

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

// Underscore from "ThisIsAName" to "this_is_a_name".
// used to format json field names
func (instance *StringsHelper) Underscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			if m != 0 {
				if !w {
					b.indent()
					w = true
				}
				b.write(m)
			}
			m = unicode.ToLower(ch)
		} else {
			if m != 0 {
				b.indent()
				b.write(m)
				m = 0
				w = false
			}
			b.write(ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}
	return string(b.r)
}

//----------------------------------------------------------------------------------------------------------------------
//	C a m e l    C a s e
//----------------------------------------------------------------------------------------------------------------------

func (instance *StringsHelper) CapitalizeAll(text string) string {
	return strings.Title(text)
}

func (instance *StringsHelper) CapitalizeFirst(text string) string {
	if len(text) > 0 {
		words := instance.Split(text, " ")
		if len(words) > 0 {
			words[0] = strings.Title(words[0])
			return instance.ConcatSep(" ", words)
		}
	}
	return text
}

//----------------------------------------------------------------------------------------------------------------------
//	q u o t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *StringsHelper) Quote(v interface{}) string {
	return strconv.Quote(Convert.ToString(v))
}

func (instance *StringsHelper) Unquote(v interface{}) (string, error) {
	return strconv.Unquote(Convert.ToString(v))
}

//----------------------------------------------------------------------------------------------------------------------
//	p a d d i n g
//----------------------------------------------------------------------------------------------------------------------

func (instance *StringsHelper) FillLeft(text string, l int, r rune) string {
	if len(text) == l {
		return text
	} else if len(text) < l {
		return fmt.Sprintf("%"+string(r)+strconv.Itoa(l)+"s", text)
	}
	return text[:l]
}

func (instance *StringsHelper) FillRight(text string, l int, r rune) string {
	if len(text) == l {
		return text
	} else if len(text) < l {
		return text + strings.Repeat(string(r), l-len(text))
	}
	return text[:l]
}

func (instance *StringsHelper) FillLeftBytes(bytes []byte, l int, r rune) []byte {
	return []byte(instance.FillLeft(string(bytes), l, r))
}

func (instance *StringsHelper) FillLeftZero(text string, l int) string {
	return instance.FillLeft(text, l, '0')
}

func (instance *StringsHelper) FillLeftBytesZero(bytes []byte, l int) []byte {
	return []byte(instance.FillLeftZero(string(bytes), l))
}

func (instance *StringsHelper) FillRightZero(text string, l int) string {
	return instance.FillRight(text, l, '0')
}

func (instance *StringsHelper) FillRightBytes(bytes []byte, l int, r rune) []byte {
	return []byte(instance.FillRight(string(bytes), l, r))
}

func (instance *StringsHelper) FillRightBytesZero(bytes []byte, l int) []byte {
	return []byte(instance.FillRight(string(bytes), l, '0'))
}

func (instance *StringsHelper) Repeat(s string, count int) string {
	if count > 0 {
		if count > 0 && len(s)*count/count != len(s) {
			panic("strings: Repeat count causes overflow")
		}

		b := make([]byte, len(s)*count)
		bp := copy(b, s)
		for bp < len(b) {
			copy(b[bp:], b[:bp])
			bp *= 2
		}
		return string(b)
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
