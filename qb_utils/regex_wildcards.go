package qb_utils

import (
	"fmt"
	"regexp"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	tags
//----------------------------------------------------------------------------------------------------------------------

type ExpressionTag struct {
	Expression    string  `json:"expression"`
	SuccessWeight float32 `json:"success-weight"` // applied to existing keywords
	FailWeight    float32 `json:"fails-weight"`   // applied for missing keywords
}

func (instance *ExpressionTag) Json() string {
	return JSON.Stringify(instance)
}

func (instance *ExpressionTag) String() string {
	return fmt.Sprintf("%s:%v:%v", instance.Expression, instance.SuccessWeight, instance.FailWeight)
}

func (instance *ExpressionTag) SuccessWeightOf(value float32) float32 {
	if instance.SuccessWeight == -1 {
		return value
	}
	return value * instance.SuccessWeight
}

func (instance *ExpressionTag) FailWeightOf(value float32) float32 {
	weight := instance.FailWeight
	if weight == -1 {
		weight = instance.SuccessWeight
	}
	if weight == -1 {
		return value
	}
	return value * weight
}

func ParseExpression(text string) (tag *ExpressionTag) {
	tag = new(ExpressionTag)
	tokens := strings.Split(text, ":")
	if len(tokens) > 0 {
		tag.Expression = tokens[0]
		tag.SuccessWeight = Convert.ToFloat32(Arrays.GetAt(tokens, 1, -1))
		tag.FailWeight = Convert.ToFloat32(Arrays.GetAt(tokens, 2, -1))
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	n l p    u t i l s
//----------------------------------------------------------------------------------------------------------------------

// WildcardScoreAll Calculate a matching score between a phrase and a check test using expressions.
// ALL expressions are evaluated.
// Failed expressions  add negative score to result
// @param [string] phrase. "hello humanity!! I'm Mario rossi"
// @param [string] expressions. All expressions to match. ["hel??0 h*", "I* * ros*"]
// 		Supported expressions:
//			["hel??0 h*", "I* * ros*"]
//			["hel??0 h*:1:1", "I* * ros*:1:0"]
func (instance *RegexHelper) WildcardScoreAll(phrase string, expressions []string) float32 {

	countWords := float32(len(Strings.Split(phrase, " ")))
	countExpressions := float32(len(expressions))

	expressionScore := float32(0)
	failScore := float32(0)
	for _, expression := range expressions {
		tag := ParseExpression(expression)
		matchedWords := instance.WildcardMatch(phrase, tag.Expression)
		if len(matchedWords) > 0 {
			// matching
			for _, matched := range matchedWords {
				countWordsMatched := float32(len(Strings.Split(matched, " ")))
				expressionScore += tag.SuccessWeightOf(countWordsMatched / countWords)
			}
		} else {
			// no matching
			failScore += tag.FailWeightOf(1 / countExpressions)
		}
	}

	return expressionScore - failScore
}

// WildcardScoreAny Calculate a matching score between a phrase and a check test using expressions.
// ALL expressions are evaluated.
// Failed expressions  do not add negative score to result
// @param [string] phrase. "hello humanity!! I'm Mario rossi"
// @param [string] expressions. All expressions to match. ["hel??0 h*", "I* * ros*"]
func (instance *RegexHelper) WildcardScoreAny(phrase string, expressions []string) float32 {
	countWords := float32(len(Strings.Split(phrase, " ")))
	// countExpressions := float32(len(expressions))

	expressionScore := float32(0)
	for _, expression := range expressions {
		tag := ParseExpression(expression)
		matchedWords := instance.WildcardMatch(phrase, tag.Expression)
		if len(matchedWords) > 0 {
			// matching
			for _, matched := range matchedWords {
				countWordsMatched := float32(len(Strings.Split(matched, " ")))
				expressionScore += tag.SuccessWeightOf(countWordsMatched / countWords)
			}
		}
	}

	return expressionScore
}

// WildcardScoreBest Calculate a matching score between a phrase and a check test using expressions.
// ALL expressions are evaluated.
// Failed expressions  do not add negative score to result.
// Return best score above all
// @param [string] phrase. "hello humanity!! I'm Mario rossi"
// @param [string] expressions. All expressions to match. ["hel??0 h*", "I* * ros*"]
func (instance *RegexHelper) WildcardScoreBest(phrase string, expressions []string) float32 {
	countWords := float32(len(Strings.Split(phrase, " ")))
	// countExpressions := float32(len(expressions))

	expressionScore := float32(0)
	for _, expression := range expressions {
		tag := ParseExpression(expression)
		matchedWords := instance.WildcardMatch(phrase, tag.Expression)
		if len(matchedWords) > 0 {
			// matching
			for _, matched := range matchedWords {
				countWordsMatched := float32(len(Strings.Split(matched, " ")))
				score := tag.SuccessWeightOf(countWordsMatched / countWords)
				if score > expressionScore {
					expressionScore = score
				}
			}
		}
	}

	return expressionScore
}

//----------------------------------------------------------------------------------------------------------------------
//	w i l d c a r d    l o o k u p
//----------------------------------------------------------------------------------------------------------------------

func (instance *RegexHelper) WildcardMatchAll(text, expression string) ([]string, [][]int) {
	exp := toRegexp(expression)
	return matchAll(text, exp)
}

func (instance *RegexHelper) WildcardMatch(text, expression string) []string {
	exp := toRegexp(expression)
	return matchString(text, exp)
}

func (instance *RegexHelper) WildcardMatchIndex(text, expression string) [][]int {
	exp := toRegexp(expression)
	return matchIndex(text, exp)
}

func (instance *RegexHelper) WildcardMatchBetween(text string, offset int, patternStart string, patternEnd string, cutset string) []string {

	expStart := toRegexp(patternStart)
	expEnd := toRegexp(patternEnd)

	return matchBetween(text, offset, expStart, expEnd, cutset)
}

// WildcardIndex Return index array of matching expression in a text starting search from offset position
// @param text string. "hello humanity!!"
// @param pattern string "hu?an*"
// @param offset int number of characters to exclude from search
// @return []int
func (instance *RegexHelper) WildcardIndex(text string, pattern string, offset int) []int {
	regex := toRegexp(pattern)
	return index(text, regex, offset)
}

// WildcardIndexLenPair Return array of pair index:word_len  of matching expression in a text
// @param text string. "hello humanity!!"
// @param pattern string "hu?an*"
// @return [][]int ex: [[12,3], [22,4]]
func (instance *RegexHelper) WildcardIndexLenPair(text string, pattern string, offset int) [][]int {
	regex := toRegexp(pattern)
	return indexLenPair(text, regex, offset)
}

// WildcardIndexArray works like WildcardIndex but cycle on an array of pattern elements
func (instance *RegexHelper) WildcardIndexArray(text string, patterns []string, offset int) [][]int {
	response := make([][]int, 0)
	for _, pattern := range patterns {
		match := instance.WildcardIndex(text, pattern, offset)
		if len(match) > 0 {
			response = append(response, match)
		}
	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	private
//----------------------------------------------------------------------------------------------------------------------

func toRegexp(wildCardExpr string) *regexp.Regexp {
	if wildCardExpr == "\n" {
		return regexp.MustCompile(wildCardExpr)
	} else {
		prefix := "(\\b"
		suffix := "\\b)"

		// replace spaces
		wildCardExpr = strings.Replace(wildCardExpr, " ", "\\W", -1)

		// escape dot
		wildCardExpr = strings.Replace(wildCardExpr, ".", "\\.", -1)

		// replace ?
		wildCardExpr = strings.Replace(wildCardExpr, "?", "(?:(?:.|\n))?", -1)

		// replace *
		wildCardExpr = strings.Replace(wildCardExpr, "*", ".*?", -1)

		r := fmt.Sprintf("%s%s%s", prefix, wildCardExpr, suffix)
		return regexp.MustCompile(r)
	}
}
