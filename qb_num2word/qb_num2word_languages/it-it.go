package qb_num2word_languages

import (
	"fmt"
	"strings"
)

func init() {
	// register the language
	Languages["it-it"] = Language{
		Name:           "Italian",
		Aliases:        []string{"it", "it-it", "it_IT", "italian"},
		Flag:           "🇮🇹",
		IntegerToWords: IntegerToItIt,
	}
}

// IntegerToItIt converts an integer to Italian words
func IntegerToItIt(input int) string {
	var italianMegas = [][]string{{"", ""}, {"mille", "mila"}, {"milione", "milioni"}, {"miliardo", "miliardi"}, {"triliardo", "triliardi"}, {"quadrillione", "quadrillioni"}, {"quintillione", "quintillioni"}, {"sextillione", "sextillioni"}, {"septillione", "septilioni"}, {"octillione", "octilioni"}, {"nonillione", "nonilioni"}, {"decillione", "decilioni"}}
	var italianUnits = []string{"", "uno", "due", "tre", "quattro", "cinque", "sei", "sette", "otto", "nove"}
	var italianTens = []string{"", "dieci", "venti", "trenta", "quaranta", "cinquanta", "sessenta", "settanta", "ottanta", "novanta"}
	var italianTeens = []string{"dieci", "undici", "dodici", "tredici", "quattordici", "quindici", "sedici", "diciasette", "dicioto", "diciannove"}

	//log.Printf("Input: %d\n", input)
	words := []string{}

	if input < 0 {
		words = append(words, "meno")
		input *= -1
	}

	// split integer in triplets
	triplets := integerToTriplets(input)
	//log.Printf("Triplets: %v\n", triplets)

	// zero is a special case
	if len(triplets) == 0 {
		return "zero"
	}

	// iterate over triplets
	for idx := len(triplets) - 1; idx >= 0; idx-- {
		triplet := triplets[idx]
		//log.Printf("Triplet: %d (idx=%d)\n", triplet, idx)

		// nothing here for empty triplet
		if triplet == 0 {
			continue
		}

		// three-digits
		hundreds := triplet / 100 % 10
		tens := triplet / 10 % 10
		units := triplet % 10
		//log.Printf("Hundreds:%d, Tens:%d, Units:%d\n", hundreds, tens, units)
		switch hundreds {
		case 0:
			break
		case 1:
			words = append(words, "cento")
			break
		default:
			words = append(words, fmt.Sprintf("%scento", italianUnits[hundreds]))
			break
		}

		if tens == 0 && units == 0 {
			goto tripletEnd
		}

		switch tens {
		case 0:
			if len(triplets) == 1 || triplet > 1 {
				// uno
				words = append(words, italianUnits[units])
			}
		case 1:
			words = append(words, italianTeens[units])
			break
		default:
			if units > 0 {
				word := fmt.Sprintf("%s%s", italianTens[tens], italianUnits[units])
				words = append(words, word)
			} else {
				words = append(words, italianTens[tens])
			}
			break
		}

	tripletEnd:
		// mega
		var mega string
		if triplet == 1 {
			mega = italianMegas[idx][0]
		} else {
			mega = italianMegas[idx][1]
		}
		if mega != "" {
			words = append(words, mega)
		}
	}

	//log.Printf("Words length: %d\n", len(words))

	return strings.Join(words, " ")
}
