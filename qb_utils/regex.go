package qb_utils

import (
	"encoding/json"
	"html"
	"regexp"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	const
//----------------------------------------------------------------------------------------------------------------------

// Regular expression patterns
const (
	DatePattern           = `(?i)(?:[0-3]?\d(?:st|nd|rd|th)?\s+(?:of\s+)?(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)|(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)\s+[0-3]?\d(?:st|nd|rd|th)?)(?:\,)?\s*(?:\d{4})?|[0-3]?\d[-\./][0-3]?\d[-\./]\d{2,4}`
	TimePattern           = `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`
	PhonePattern          = `(?:(?:\+?\d{1,3}[-.\s*]?)?(?:\(?\d{3}\)?[-.\s*]?)?\d{3}[-.\s*]?\d{4,6})|(?:(?:(?:\(\+?\d{2}\))|(?:\+?\d{2}))\s*\d{2}\s*\d{3}\s*\d{4})`
	PhonesWithExtsPattern = `(?i)(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?(?:[2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?(?:[0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(?:\d+)?)`
	LinkPattern           = `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`
	LinkPatternStrict     = `(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`
	EmailPattern          = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	IPv4Pattern           = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	IPv6Pattern           = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	IPPattern             = IPv4Pattern + `|` + IPv6Pattern
	NotKnownPortPattern   = `6[0-5]{2}[0-3][0-5]|[1-5][\d]{4}|[2-9][\d]{3}|1[1-9][\d]{2}|10[3-9][\d]|102[4-9]`
	//PricePattern          = `[$]\s?[+-]?[0-9]{1,3}(?:(?:,?[0-9]{3}))*(?:\.[0-9]{1,2})?`
	PricePattern          = `[$€]?\s?[+-]?[0-9]{1,3}(?:(?:(,|\.)?[0-9]{3}))*(?:(\.|,)[0-9]{1,3})?\s?[$€]?`
	HexColorPattern       = `(?:#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3}))`
	CreditCardPattern     = `(?:(?:(?:\d{4}[- ]?){3}\d{4}|\d{15,16}))`
	VISACreditCardPattern = `4\d{3}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	MCCreditCardPattern   = `5[1-5]\d{2}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	BtcAddressPattern     = `[13][a-km-zA-HJ-NP-Z1-9]{25,34}`
	StreetAddressPattern  = `\d{1,4} [\w\s]{1,20}(?:street|st|avenue|ave|road|rd|highway|hwy|square|sq|trail|trl|drive|dr|court|ct|park|parkway|pkwy|circle|cir|boulevard|blvd)\W?`
	ZipCodePattern        = `\b\d{5}(?:[-\s]\d{4})?\b`
	PoBoxPattern          = `(?i)P\.? ?O\.? Box \d+`
	SSNPattern            = `(?:\d{3}-\d{2}-\d{4})`
	MD5HexPattern         = `[0-9a-fA-F]{32}`
	SHA1HexPattern        = `[0-9a-fA-F]{40}`
	SHA256HexPattern      = `[0-9a-fA-F]{64}`
	GUIDPattern           = `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`
	ISBN13Pattern         = `(?:[\d]-?){12}[\dxX]`
	ISBN10Pattern         = `(?:[\d]-?){9}[\dxX]`
	MACAddressPattern     = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
	IBANPattern           = `[A-Z]{2}\d{2}[A-Z0-9]{4}\d{7}([A-Z\d]?){0,16}`
	GitRepoPattern        = `((git|ssh|http(s)?)|(git@[\w\.]+))(:(\/\/)?)([\w\.@\:/\-~]+)(\.git)(\/)?`
	NumbersPattern        = `[-+]?(?:[0-9]+(\.|,))*[0-9]+(?:\.[0-9]+)?`
	BracesPattern         = `\{{(.*?)\}}` // `{{\s*[\w\.]+\s*}}`
)

var (
	SQLBlackList = []string{"alter", "ALTER", "EXEC", "exec", "EXECUTE", "execute", "create", "CREATE",
		"delete", "DELETE", "DROP", "drop", "INSERT", "insert", "UPDATE", "update", "UNION", "union",
		"--", "\"", "'",
	}
)

// Compiled regular expressions
var (
	DateRegex           = regexp.MustCompile(DatePattern)
	TimeRegex           = regexp.MustCompile(TimePattern)
	PhoneRegex          = regexp.MustCompile(PhonePattern)
	PhonesWithExtsRegex = regexp.MustCompile(PhonesWithExtsPattern)
	LinkRegex           = regexp.MustCompile(LinkPattern)
	UrlRegex            = regexp.MustCompile(LinkPatternStrict)
	EmailRegex          = regexp.MustCompile(EmailPattern)
	IPv4Regex           = regexp.MustCompile(IPv4Pattern)
	IPv6Regex           = regexp.MustCompile(IPv6Pattern)
	IPRegex             = regexp.MustCompile(IPPattern)
	NotKnownPortRegex   = regexp.MustCompile(NotKnownPortPattern)
	PriceRegex          = regexp.MustCompile(PricePattern)
	HexColorRegex       = regexp.MustCompile(HexColorPattern)
	CreditCardRegex     = regexp.MustCompile(CreditCardPattern)
	BtcAddressRegex     = regexp.MustCompile(BtcAddressPattern)
	StreetAddressRegex  = regexp.MustCompile(StreetAddressPattern)
	ZipCodeRegex        = regexp.MustCompile(ZipCodePattern)
	PoBoxRegex          = regexp.MustCompile(PoBoxPattern)
	SSNRegex            = regexp.MustCompile(SSNPattern)
	MD5HexRegex         = regexp.MustCompile(MD5HexPattern)
	SHA1HexRegex        = regexp.MustCompile(SHA1HexPattern)
	SHA256HexRegex      = regexp.MustCompile(SHA256HexPattern)
	GUIDRegex           = regexp.MustCompile(GUIDPattern)
	ISBN13Regex         = regexp.MustCompile(ISBN13Pattern)
	ISBN10Regex         = regexp.MustCompile(ISBN10Pattern)
	VISACreditCardRegex = regexp.MustCompile(VISACreditCardPattern)
	MCCreditCardRegex   = regexp.MustCompile(MCCreditCardPattern)
	MACAddressRegex     = regexp.MustCompile(MACAddressPattern)
	IBANRegex           = regexp.MustCompile(IBANPattern)
	GitRepoRegex        = regexp.MustCompile(GitRepoPattern)
	NumbersRegex        = regexp.MustCompile(NumbersPattern)
	BracesRegex         = regexp.MustCompile(BracesPattern)
)

//----------------------------------------------------------------------------------------------------------------------
//	helper
//----------------------------------------------------------------------------------------------------------------------

type RegexHelper struct {
}

var Regex *RegexHelper

func init() {
	Regex = new(RegexHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x t r a c t o r
//----------------------------------------------------------------------------------------------------------------------

// Date finds all date strings
func (instance *RegexHelper) Date(text string) []string {
	return matchString(text, DateRegex)
}

// Time finds all time strings
func (instance *RegexHelper) Time(text string) []string {
	return matchString(text, TimeRegex)
}

// Phones finds all phone numbers
func (instance *RegexHelper) Phones(text string) []string {
	return matchString(text, PhoneRegex)
}

// PhonesWithExts finds all phone numbers with ext
func (instance *RegexHelper) PhonesWithExts(text string) []string {
	return matchString(text, PhonesWithExtsRegex)
}

// Links finds all link strings
func (instance *RegexHelper) Links(text string) []string {
	return matchString(text, LinkRegex)
}

// Urls finds all URL strings
func (instance *RegexHelper) Urls(text string) []string {
	return matchString(text, UrlRegex)
}

// Emails finds all email strings
func (instance *RegexHelper) Emails(text string) []string {
	return matchString(text, EmailRegex)
}

// IPv4s finds all IPv4 addresses
func (instance *RegexHelper) IPv4s(text string) []string {
	return matchString(text, IPv4Regex)
}

// IPv6s finds all IPv6 addresses
func (instance *RegexHelper) IPv6s(text string) []string {
	return matchString(text, IPv6Regex)
}

// IPs finds all IP addresses (both IPv4 and IPv6)
func (instance *RegexHelper) IPs(text string) []string {
	return matchString(text, IPRegex)
}

// NotKnownPorts finds all not-known port numbers
func (instance *RegexHelper) NotKnownPorts(text string) []string {
	return matchString(text, NotKnownPortRegex)
}

// Prices finds all price strings
func (instance *RegexHelper) Prices(text string) []string {
	array := matchString(text, PriceRegex)
	Strings.TrimSpaces(array)
	return array
}

// HexColors finds all hex color values
func (instance *RegexHelper) HexColors(text string) []string {
	return matchString(text, HexColorRegex)
}

// CreditCards finds all credit card numbers
func (instance *RegexHelper) CreditCards(text string) []string {
	return matchString(text, CreditCardRegex)
}

// BtcAddresses finds all bitcoin addresses
func (instance *RegexHelper) BtcAddresses(text string) []string {
	return matchString(text, BtcAddressRegex)
}

// StreetAddresses finds all street addresses
func (instance *RegexHelper) StreetAddresses(text string) []string {
	return matchString(text, StreetAddressRegex)
}

// ZipCodes finds all zip codes
func (instance *RegexHelper) ZipCodes(text string) []string {
	return matchString(text, ZipCodeRegex)
}

// PoBoxes finds all po-box strings
func (instance *RegexHelper) PoBoxes(text string) []string {
	return matchString(text, PoBoxRegex)
}

// SSNs finds all SSN strings
func (instance *RegexHelper) SSNs(text string) []string {
	return matchString(text, SSNRegex)
}

// MD5Hexes finds all MD5 hex strings
func (instance *RegexHelper) MD5Hexes(text string) []string {
	return matchString(text, MD5HexRegex)
}

// SHA1Hexes finds all SHA1 hex strings
func (instance *RegexHelper) SHA1Hexes(text string) []string {
	return matchString(text, SHA1HexRegex)
}

// SHA256Hexes finds all SHA256 hex strings
func (instance *RegexHelper) SHA256Hexes(text string) []string {
	return matchString(text, SHA256HexRegex)
}

// GUIDs finds all GUID strings
func (instance *RegexHelper) GUIDs(text string) []string {
	return matchString(text, GUIDRegex)
}

// ISBN13s finds all ISBN13 strings
func (instance *RegexHelper) ISBN13s(text string) []string {
	return matchString(text, ISBN13Regex)
}

// ISBN10s finds all ISBN10 strings
func (instance *RegexHelper) ISBN10s(text string) []string {
	return matchString(text, ISBN10Regex)
}

// VISACreditCards finds all VISA credit card numbers
func (instance *RegexHelper) VISACreditCards(text string) []string {
	return matchString(text, VISACreditCardRegex)
}

// MCCreditCards finds all MasterCard credit card numbers
func (instance *RegexHelper) MCCreditCards(text string) []string {
	return matchString(text, MCCreditCardRegex)
}

// MACAddresses finds all MAC addresses
func (instance *RegexHelper) MACAddresses(text string) []string {
	return matchString(text, MACAddressRegex)
}

// IBANs finds all IBAN strings
func (instance *RegexHelper) IBANs(text string) []string {
	return matchString(text, IBANRegex)
}

// GitRepos finds all git repository addresses which have protocol prefix
func (instance *RegexHelper) GitRepos(text string) []string {
	return matchString(text, GitRepoRegex)
}

func (instance *RegexHelper) Numbers(text string) []string {
	return matchString(text, NumbersRegex)
}

func (instance *RegexHelper) BetweenBraces(text string) [][]string {
	return BracesRegex.FindAllStringSubmatch(text, -1)
}

func (instance *RegexHelper) TagsBetweenBraces(text string) []string {
	return BracesRegex.FindAllString(text, -1)
}

func (instance *RegexHelper) TextBetweenBraces(text string) []string {
	tags := BracesRegex.FindAllStringSubmatch(text, -1)
	response := make([]string, 0)
	for _, v := range tags {
		response = append(response, v[1])
	}
	return response
}

func (instance *RegexHelper) TagsBetweenStrings(text string, prefix, suffix string) []string {
	return tagsBetweenStrings(text, prefix, suffix, false)
}

func (instance *RegexHelper) TagsBetweenTrimStrings(text string, prefix, suffix string) []string {
	return tagsBetweenStrings(text, prefix, suffix, true)
}

func (instance *RegexHelper) TextBetweenStrings(text string, prefix string, suffix interface{}) []string {
	tt := strings.ReplaceAll(text, prefix, "{{")

	if arr, ok := suffix.([]string); ok {
		for _, str := range arr {
			tt = replaceAll(tt, str)
		}
	} else if str, ok := suffix.(string); ok {
		tt = replaceAll(tt, str)
	}

	tags := BracesRegex.FindAllStringSubmatch(tt, -1)
	response := make([]string, 0)
	for _, v := range tags {
		response = append(response, v[1])
	}
	return response
}

// GetParamNames return unique @param names
func (instance *RegexHelper) GetParamNames(statement, prefix string, suffix interface{}) []string {
	response := make([]string, 0)
	statement = strings.ReplaceAll(statement, ";", " ;")
	statement += " "
	params := instance.TextBetweenStrings(statement, prefix, suffix)
	for _, param := range params {
		// purge name from comma or other invalid delimiters
		param = strings.TrimRight(param, ",.;:\n\r")
		if Arrays.IndexOf(param, response) == -1 {
			response = append(response, param)
		}
	}
	return response
}

func (instance *RegexHelper) GetParamNamesBraces(statement string) []string {
	return instance.GetParamNames(statement, "{{", "}}")
}

func (instance *RegexHelper) GetParamNamesAt(statement string) []string {
	return instance.GetParamNames(statement, "@", " ")
}

//----------------------------------------------------------------------------------------------------------------------
//	v a l i d a t i o n
//----------------------------------------------------------------------------------------------------------------------

func (instance *RegexHelper) IsValidEmail(text string) bool {
	return len(instance.Emails(text)) == 1
}

func (instance *RegexHelper) IsValidJsonObject(text string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(text), &js) == nil
}

func (instance *RegexHelper) IsValidJsonArray(text string) bool {
	var js []map[string]interface{}
	return json.Unmarshal([]byte(text), &js) == nil
}

func (instance *RegexHelper) IsHTML(text string) bool {
	if strings.Index(text, "<") > -1 && (strings.Index(text, "</") > -1 || strings.Index(text, "/>") > -1) {
		return true
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p    l o o k u p
//----------------------------------------------------------------------------------------------------------------------

func (instance *RegexHelper) MatchAll(text, expression string) ([]string, [][]int) {
	exp := regexp.MustCompile(expression)
	return matchAll(text, exp)
}

func (instance *RegexHelper) Match(text, expression string) []string {
	exp := regexp.MustCompile(expression)
	return matchString(text, exp)
}

func (instance *RegexHelper) MatchIndex(text, expression string) [][]int {
	exp := regexp.MustCompile(expression)
	return matchIndex(text, exp)
}

func (instance *RegexHelper) MatchBetween(text string, offset int, patternStart string, patternEnd string, cutset string) []string {
	expStart := regexp.MustCompile(patternStart)
	expEnd := regexp.MustCompile(patternEnd)

	return matchBetween(text, offset, expStart, expEnd, cutset)
}

func (instance *RegexHelper) Index(text string, pattern string, offset int) []int {
	regex := regexp.MustCompile(pattern)
	return index(text, regex, offset)
}

func (instance *RegexHelper) IndexLenPair(text string, pattern string, offset int) [][]int {
	regex := regexp.MustCompile(pattern)
	return indexLenPair(text, regex, offset)
}

//----------------------------------------------------------------------------------------------------------------------
//	s a n i t i z e
//----------------------------------------------------------------------------------------------------------------------

func (instance *RegexHelper) SanitizeHTML(text string) (clean string) {
	clean = html.EscapeString(text)
	return
}

func (instance *RegexHelper) SanitizeSQL(text string) (clean string) {
	clean = text
	for _, s := range SQLBlackList {
		clean = strings.ReplaceAll(clean, s, "")
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func tagsBetweenStrings(text string, prefix, suffix string, trimSpaces bool) []string {
	tt := strings.ReplaceAll(text, prefix, "{{")
	if suffix == " " {
		tt = strings.ReplaceAll(tt, suffix, "}} ")
		tt = strings.ReplaceAll(tt, "\n", "}} ")
		tt = strings.ReplaceAll(tt, "\t", "}} ")
	} else {
		tt = strings.ReplaceAll(tt, suffix, "}}")
	}

	tags := BracesRegex.FindAllString(tt, -1)
	response := make([]string, 0)
	for _, v := range tags {
		v = strings.ReplaceAll(strings.ReplaceAll(v, "}}", suffix), "{{", prefix)
		if trimSpaces {
			v = strings.TrimSpace(v)
		}
		response = append(response, v)
	}
	return response
}

func matchAll(text string, regex *regexp.Regexp) ([]string, [][]int) {
	parsed := regex.FindAllString(text, -1)
	index := regex.FindAllStringIndex(text, -1)
	return parsed, index
}

func matchString(text string, regex *regexp.Regexp) []string {
	return regex.FindAllString(text, -1)
}

func matchIndex(text string, regex *regexp.Regexp) [][]int {
	return regex.FindAllStringIndex(text, -1)
}

func matchBetween(text string, offset int, patternStart *regexp.Regexp, patternEnd *regexp.Regexp, cutset string) []string {
	text = Strings.Sub(text, offset, len(text))
	response := make([]string, 0)

	indexesStart := matchIndex(text, patternStart) // [][]int
	for _, indexStart := range indexesStart {
		is := indexStart[0]
		ie := indexStart[1] // end of first pattern
		if is < ie {
			sub := Strings.Sub(text, ie, len(text))
			indexesEnd := matchIndex(sub, patternEnd) // [][]int
			if len(indexesEnd) > 0 {
				indexEnd := indexesEnd[0][0]
				sub = Strings.Sub(sub, 0, indexEnd)
				if len(cutset) > 0 {
					sub = strings.Trim(sub, cutset)
				}
				response = append(response, sub)
			} else {
				if len(cutset) > 0 {
					sub = strings.Trim(sub, cutset)
				}
				response = append(response, sub)
			}
		}
	}

	return response
}

func index(text string, regex *regexp.Regexp, offset int) []int {
	var response []int
	if nil != regex && len(text) > 0 {
		if offset < 0 {
			offset = 0
		}

		// shrink text starting from offset
		text = Strings.Sub(text, offset, len(text))

		// get regexp match
		indexes := matchIndex(text, regex)

		if len(indexes) > 0 {
			for _, index := range indexes {
				response = append(response, index[0]+offset)
			}
		}
	}
	return response
}

func indexLenPair(text string, regex *regexp.Regexp, offset int) [][]int {
	var response [][]int
	if nil != regex && len(text) > 0 {
		if offset < 0 {
			offset = 0
		}

		// shrink text starting from offset
		text = Strings.Sub(text, offset, len(text))

		// get regexp match
		indexes := matchIndex(text, regex)

		if len(indexes) > 0 {
			for _, index := range indexes {
				pair := make([]int, 2)
				pair[0] = index[0] + offset
				pair[1] = index[1] - index[0]
				response = append(response, pair)
			}
		}
	}
	return response
}

func replaceAll(source, str string) (tt string) {
	if str == " " {
		tt = strings.ReplaceAll(source, str, "}} ")
		tt = strings.ReplaceAll(tt, "\n", "}} ")
		tt = strings.ReplaceAll(tt, "\t", "}} ")
	} else {
		tt = strings.ReplaceAll(source, str, "}}")
	}
	return
}
