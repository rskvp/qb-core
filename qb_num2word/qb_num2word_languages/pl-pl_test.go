package qb_num2word_languages

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerToPlPl(t *testing.T) {
	t.Parallel()

	tests := map[int]string{
		-1:            "minus jeden",
		0:             "zero",
		1:             "jeden",
		2:             "dwa",
		3:             "trzy",
		4:             "cztery",
		5:             "pięć",
		6:             "sześć",
		7:             "siedem",
		8:             "osiem",
		9:             "dziewięć",
		10:            "dziesięć",
		11:            "jedenaście",
		12:            "dwanaście",
		13:            "trzynaście",
		14:            "czternaście",
		15:            "piętnaście",
		16:            "szesnaście",
		17:            "siedemnaście",
		18:            "osiemnaście",
		19:            "dziewiętnaście",
		20:            "dwadzieścia",
		21:            "dwadzieścia jeden",
		22:            "dwadzieścia dwa",
		23:            "dwadzieścia trzy",
		24:            "dwadzieścia cztery",
		25:            "dwadzieścia pięć",
		26:            "dwadzieścia sześć",
		27:            "dwadzieścia siedem",
		28:            "dwadzieścia osiem",
		29:            "dwadzieścia dziewięć",
		30:            "trzydzieści",
		31:            "trzydzieści jeden",
		32:            "trzydzieści dwa",
		39:            "trzydzieści dziewięć",
		42:            "czterdzieści dwa",
		80:            "osiemdziesiąt",
		90:            "dziewięćdziesiąt",
		99:            "dziewięćdziesiąt dziewięć",
		100:           "sto",
		101:           "sto jeden",
		111:           "sto jedenaście",
		120:           "sto dwadzieścia",
		121:           "sto dwadzieścia jeden",
		200:           "dwieście",
		300:           "trzysta",
		400:           "czterysta",
		500:           "pięćset",
		600:           "sześćset",
		700:           "siedemset",
		800:           "osiemset",
		900:           "dziewięćset",
		909:           "dziewięćset dziewięć",
		919:           "dziewięćset dziewiętnaście",
		990:           "dziewięćset dziewięćdziesiąt",
		999:           "dziewięćset dziewięćdziesiąt dziewięć",
		1000:          "jeden tysiąc",
		1337:          "jeden tysiąc trzysta trzydzieści siedem",
		2000:          "dwa tysiące",
		4000:          "cztery tysiące",
		5000:          "pięć tysięcy",
		11000:         "jedenaście tysięcy",
		21000:         "dwadzieścia jeden tysięcy",
		28000:         "dwadzieścia osiem tysięcy",
		31000:         "trzydzieści jeden tysięcy",
		32000:         "trzydzieści dwa tysiące",
		39000:         "trzydzieści dziewięć tysięcy",
		42000:         "czterdzieści dwa tysiące",
		999000:        "dziewięćset dziewięćdziesiąt dziewięć tysięcy",
		999999:        "dziewięćset dziewięćdziesiąt dziewięć tysięcy dziewięćset dziewięćdziesiąt dziewięć",
		1000000:       "jeden milion",
		2000000:       "dwa miliony",
		4000000:       "cztery miliony",
		5000000:       "pięć milionów",
		100100100:     "sto milionów sto tysięcy sto",
		500500500:     "pięćset milionów pięćset tysięcy pięćset",
		606606606:     "sześćset sześć milionów sześćset sześć tysięcy sześćset sześć",
		999000000:     "dziewięćset dziewięćdziesiąt dziewięć milionów",
		999000999:     "dziewięćset dziewięćdziesiąt dziewięć milionów dziewięćset dziewięćdziesiąt dziewięć",
		999999000:     "dziewięćset dziewięćdziesiąt dziewięć milionów dziewięćset dziewięćdziesiąt dziewięć tysięcy",
		999999999:     "dziewięćset dziewięćdziesiąt dziewięć milionów dziewięćset dziewięćdziesiąt dziewięć tysięcy dziewięćset dziewięćdziesiąt dziewięć",
		1174315110:    "jeden miliard sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziesięć",
		1174315119:    "jeden miliard sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziewiętnaście",
		15174315119:   "piętnaście miliardów sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziewiętnaście",
		35174315119:   "trzydzieści pięć miliardów sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziewiętnaście",
		935174315119:  "dziewięćset trzydzieści pięć miliardów sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziewiętnaście",
		2935174315119: "dwa biliony dziewięćset trzydzieści pięć miliardów sto siedemdziesiąt cztery miliony trzysta piętnaście tysięcy sto dziewiętnaście",
	}

	for input, expectedOutput := range tests {
		name := fmt.Sprintf("%d", input)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expectedOutput, IntegerToPlPl(input))
		})
	}
}
