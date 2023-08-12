package qb_num2word_languages

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleIntegerToTrTr() {
	fmt.Println(IntegerToTrTr(42))
	// Output: kırk iki
}

func TestIntegerToTrTr(t *testing.T) {
	t.Parallel()

	tests := map[int]string{
		-1:            "eksi bir",
		0:             "sıfır",
		1:             "bir",
		9:             "dokuz",
		10:            "on",
		11:            "on bir",
		19:            "on dokuz",
		20:            "yirmi",
		21:            "yirmi bir",
		80:            "seksen",
		90:            "doksan",
		99:            "doksan dokuz",
		100:           "yüz",
		101:           "yüz bir",
		111:           "yüz on bir",
		120:           "yüz yirmi",
		121:           "yüz yirmi bir",
		900:           "dokuz yüz",
		909:           "dokuz yüz dokuz",
		919:           "dokuz yüz on dokuz",
		990:           "dokuz yüz doksan",
		999:           "dokuz yüz doksan dokuz",
		1000:          "bin",
		2000:          "iki bin",
		4000:          "dört bin",
		5000:          "beş bin",
		11000:         "on bir bin",
		21000:         "yirmi bir bin",
		999000:        "dokuz yüz doksan dokuz bin",
		999999:        "dokuz yüz doksan dokuz bin dokuz yüz doksan dokuz",
		1000000:       "bir milyon",
		2000000:       "iki milyon",
		4000000:       "dört milyon",
		5000000:       "beş milyon",
		100100100:     "yüz milyon yüz bin yüz",
		500500500:     "beş yüz milyon beş yüz bin beş yüz",
		606606606:     "altı yüz altı milyon altı yüz altı bin altı yüz altı",
		999000000:     "dokuz yüz doksan dokuz milyon",
		999000999:     "dokuz yüz doksan dokuz milyon dokuz yüz doksan dokuz",
		999999000:     "dokuz yüz doksan dokuz milyon dokuz yüz doksan dokuz bin",
		999999999:     "dokuz yüz doksan dokuz milyon dokuz yüz doksan dokuz bin dokuz yüz doksan dokuz",
		1174315110:    "bir milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on",
		1174315119:    "bir milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on dokuz",
		15174315119:   "on beş milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on dokuz",
		35174315119:   "otuz beş milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on dokuz",
		935174315119:  "dokuz yüz otuz beş milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on dokuz",
		2935174315119: "iki trilyon dokuz yüz otuz beş milyar yüz yetmiş dört milyon üç yüz on beş bin yüz on dokuz",
	}

	for input, expectedOutput := range tests {
		name := fmt.Sprintf("%d", input)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expectedOutput, IntegerToTrTr(input))
		})
	}
}
