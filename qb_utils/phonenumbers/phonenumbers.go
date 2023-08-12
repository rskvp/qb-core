package phonenumbers

import (
	"regexp"
	"strings"
)

type PhoneNumberHelper struct {
}

var PhoneNumber *PhoneNumberHelper

func init() {
	PhoneNumber = new(PhoneNumberHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Parse mobile number by country
func (instance *PhoneNumberHelper) Parse(number string, country string) string {
	return instance.parseInternal(number, country, false)
}

// ParseWithLandLine is Parse mobile and land line number by country
func (instance *PhoneNumberHelper) ParseWithLandLine(number string, country string) string {
	return instance.parseInternal(number, country, true)
}

// GetISO3166ByNumber ...
func (instance *PhoneNumberHelper) GetISO3166ByNumber(number string, withLandLine bool) ISO3166 {
	iso3166 := ISO3166{}
	for _, i := range GetISO3166() {
		r := regexp.MustCompile(`^` + i.CountryCode)
		for _, l := range i.PhoneNumberLengths {
			if r.MatchString(number) && len(number) == len(i.CountryCode)+l {
				// Check match with mobile codes
				for _, w := range i.MobileBeginWith {
					r := regexp.MustCompile(`^` + i.CountryCode + w)
					if r.MatchString(number) {
						// Match by mobile codes
						return i
					}
				}

				// Match by country code only for land line numbers only
				if withLandLine == true {
					return i
				}
			}
		}
	}
	return iso3166
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *PhoneNumberHelper)  parseInternal(number string, country string, landLineInclude bool) string {
	number = strings.Replace(number, " ", "", -1)
	country = strings.Replace(country, " ", "", -1)
	plusSign := false
	if strings.HasPrefix(number, "+") {
		plusSign = true
	}

	// remove any non-digit character, included the +
	number = regexp.MustCompile(`\D`).ReplaceAllString(number, "")

	iso3166 := instance.GetISO3166ByCountry(country)

	if instance.indexOfString(iso3166.Alpha3, []string{"GAB", "CIV", "COG"}) == -1 {
		r := regexp.MustCompile(`^0+`)
		number = r.ReplaceAllString(number, "")
	}
	r := regexp.MustCompile(`^89`)
	if iso3166.Alpha3 == "RUS" && len(number) == 11 && r.MatchString(number) == true {
		r := regexp.MustCompile(`^8+`)
		number = r.ReplaceAllString(number, "")
	}
	if plusSign {
		iso3166 = instance.GetISO3166ByNumber(number, landLineInclude)
	} else {
		if instance.indexOfInt(len(number), iso3166.PhoneNumberLengths) != -1 {
			number = iso3166.CountryCode + number
		}
	}
	if instance.validatePhoneISO3166(number, iso3166, landLineInclude) {
		return number
	}
	return ""
}

func (instance *PhoneNumberHelper) GetISO3166ByCountry(country string) ISO3166 {
	iso3166 := ISO3166{}
	uppperCaseCountry := strings.ToUpper(country)
	switch len(country) {
	case 0:
		iso3166 = GetISO3166()[0]
		break
	case 2:
		for _, i := range GetISO3166() {
			if i.Alpha2 == uppperCaseCountry {
				iso3166 = i
				break
			}
		}
		break
	case 3:
		for _, i := range GetISO3166() {
			if i.Alpha3 == uppperCaseCountry {
				iso3166 = i
				break
			}
		}
		break
	default:
		for _, i := range GetISO3166() {
			if strings.ToUpper(i.CountryName) == uppperCaseCountry {
				iso3166 = i
				break
			}
		}
		break
	}
	return iso3166
}

func (instance *PhoneNumberHelper) validatePhoneISO3166(number string, iso3166 ISO3166, withLandLine bool) bool {
	if len(iso3166.PhoneNumberLengths) == 0 {
		return false
	}

	if withLandLine {
		r := regexp.MustCompile(`^` + iso3166.CountryCode)
		for _, l := range iso3166.PhoneNumberLengths {
			if r.MatchString(number) && len(number) == len(iso3166.CountryCode)+l {
				return true
			}
		}
		return false
	}

	r := regexp.MustCompile(`^` + iso3166.CountryCode)
	number = r.ReplaceAllString(number, "")
	for _, l := range iso3166.PhoneNumberLengths {
		if l == len(number) {
			for _, w := range iso3166.MobileBeginWith {
				r := regexp.MustCompile(`^` + w)
				if r.MatchString(number) == true {
					return true
				}
			}
		}
	}
	return false
}

func (instance *PhoneNumberHelper) indexOfString(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func (instance *PhoneNumberHelper) indexOfInt(word int, data []int) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

