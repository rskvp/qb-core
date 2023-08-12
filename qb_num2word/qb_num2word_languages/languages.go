package qb_num2word_languages

import (
	"fmt"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Language struct {
	Name    string
	Aliases []string
	Flag    string
	IntegerToWords func(int) string
}

type LanguageList map[string]Language

var Languages = LanguageList{}

func (lang Language) HelpText() string {
	output := fmt.Sprintf("%s (%s) %s", lang.Name, strings.Join(lang.Aliases, ", "), lang.Flag)
	if lang.Name == Languages.Default().Name {
		output += "  *default*"
	}
	return output
}

func (langs LanguageList) Default() Language {
	return langs["en-us"]
}

func (langs LanguageList) Lookup(key string) *Language {
	for _, lang := range langs {
		for _, alias := range lang.Aliases {
			if alias == key {
				return &lang
			}
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	pr i v a t e
//----------------------------------------------------------------------------------------------------------------------

func integerToTriplets(number int) []int {
	triplets := []int{}

	for number > 0 {
		triplets = append(triplets, number%1000)
		number = number / 1000
	}

	return triplets
}


