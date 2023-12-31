package qb_num2word_languages

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleIntegerToJaJp() {
	fmt.Println(IntegerToJaJp(42))
	// Output: 四十二
}

func TestIntegerToJaJp(t *testing.T) {
	t.Parallel()

	tests := map[int]string{
		-1:            "マイナス一",
		0:             "零",
		1:             "一",
		9:             "九",
		10:            "十",
		11:            "十一",
		19:            "十九",
		20:            "二十",
		21:            "二十一",
		80:            "八十",
		90:            "九十",
		99:            "九十九",
		100:           "百",
		101:           "百一",
		111:           "百十一",
		120:           "百二十",
		121:           "百二十一",
		900:           "九百",
		909:           "九百九",
		919:           "九百十九",
		990:           "九百九十",
		999:           "九百九十九",
		1000:          "千",
		2000:          "二千",
		4000:          "四千",
		5000:          "五千",
		11000:         "一万千",
		21000:         "二万千",
		999000:        "九十九万九千",
		999999:        "九十九万九千九百九十九",
		1000000:       "百万",
		2000000:       "二百万",
		4000000:       "四百万",
		5000000:       "五百万",
		100100100:     "一億十万百",
		500500500:     "五億五十万五百",
		606606606:     "六億六百六十万六千六百六",
		999000000:     "九億九千九百万",
		999000999:     "九億九千九百万九百九十九",
		999990009:     "九億九千九百九十九万九",
		999999999:     "九億九千九百九十九万九千九百九十九",
		1174315110:    "十一億七千四百三十一万五千百十",
		1174315119:    "十一億七千四百三十一万五千百十九",
		15174315119:   "百五十一億七千四百三十一万五千百十九",
		35174315119:   "三百五十一億七千四百三十一万五千百十九",
		935174315119:  "九千三百五十一億七千四百三十一万五千百十九",
		2935174315119: "二兆九千三百五十一億七千四百三十一万五千百十九",
	}

	for input, expectedOutput := range tests {
		name := fmt.Sprintf("%d", input)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expectedOutput, IntegerToJaJp(input))
		})
	}
}
