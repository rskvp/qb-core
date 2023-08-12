package genusers

import (
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_utils"
)

type GenUsersHelper struct {
}

var GenUsers *GenUsersHelper

func init() {
	GenUsers = new(GenUsersHelper)
}

func (instance *GenUsersHelper) NewEngine(root string) *GenUsersEngine {
	response := new(GenUsersEngine)

	response.init(root)

	return response
}

type GenUserItem struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	FullName    string `json:"fullname"`
	Gender      string `json:"gender"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (instance *GenUserItem) String() string {
	return qb_utils.JSON.Stringify(instance)
}

func (instance *GenUserItem) Map(fields ...string) map[string]interface{} {
	m := qb_utils.Convert.ToMap(instance)
	if len(fields) == 0 ||
		(len(fields) == 1 && fields[0] == "*" || len(fields) == 1 && fields[0] == "" || len(fields) == 1 && fields[0] == "all") {
		return m
	}
	mm := make(map[string]interface{})
	for _, field := range fields {
		if v, b := m[field]; b {
			mm[field] = v
		}
	}
	return mm
}

//----------------------------------------------------------------------------------------------------------------------
//	GenUsersEngine
//----------------------------------------------------------------------------------------------------------------------

type GenUsersEngine struct {
	root             string
	fileNames        string
	fileSurnames     string
	fileCountryCodes string
	maxNames         int
	maxSurnames      int
	maxCountryCodes  int
	err              error
	options          *qb_utils.CsvOptions
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *GenUsersEngine) HasError() bool {
	return nil != instance.err
}

func (instance *GenUsersEngine) Generate(num int) ([]*GenUserItem, error) {
	response := make([]*GenUserItem, 0)
	if !instance.HasError() {
		for i := 0; i < num; i++ {
			idn := int(qb_rnd.Rnd.Between(0, int64(instance.maxNames)))
			name, gender, en := instance.readName(idn)
			if nil != en {
				return response, en
			}
			ids := int(qb_rnd.Rnd.Between(0, int64(instance.maxSurnames)))
			surname, es := instance.readSurname(ids)
			if nil != es {
				return response, es
			}
			fullname := fmt.Sprintf("%s %s", name, surname)
			username := fmt.Sprintf("%s.%s", strings.ToLower(name), strings.ToLower(surname))
			email := fmt.Sprintf("%s@%s", username, "text.gq")
			idc := int(qb_rnd.Rnd.Between(0, int64(instance.maxCountryCodes)))
			country, code, prefix, ec := instance.readCountryCode(idc)
			if nil != ec {
				return response, ec
			}
			mobile := fmt.Sprintf("%s%s", prefix, qb_rnd.Rnd.RndDigits(10))
			// password
			password := qb_rnd.Rnd.RndChars(8)
			// add user
			response = append(response, &GenUserItem{
				Name:        name,
				Surname:     surname,
				FullName:    fullname,
				Gender:      gender,
				Email:       email,
				Mobile:      mobile,
				Country:     country,
				CountryCode: code,
				Username:    username,
				Password:    password,
			})
		}
	}
	return response, instance.err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *GenUsersEngine) init(root string) {
	// work in application workspace
	if len(root) > 0 {
		instance.root = qb_utils.Paths.Concat(qb_utils.Paths.Absolute(root), "./gen_users")
	} else {
		instance.root = qb_utils.Paths.WorkspacePath("./gen_users")
	}
	instance.fileNames = qb_utils.Paths.Concat(instance.root, "names.csv")
	instance.fileSurnames = qb_utils.Paths.Concat(instance.root, "surnames.csv")
	instance.fileCountryCodes = qb_utils.Paths.Concat(instance.root, "country_codes.csv")
	_ = qb_utils.Paths.Mkdir(instance.fileNames)

	// check exists
	if b, _ := qb_utils.Paths.Exists(instance.fileNames); !b {
		// download
		_, errs := DownloadTemplates(qb_utils.Paths.Dir(instance.fileNames))
		if len(errs) > 0 {
			instance.err = qb_utils.Errors.Prefix(errs[0], "Download Error: ")
		}
	}

	// count
	if b, _ := qb_utils.Paths.Exists(instance.fileNames); b {
		_ = qb_utils.IO.ScanTextFromFile(instance.fileNames, func(_ string) bool {
			instance.maxNames++
			return false
		})
	}
	if b, _ := qb_utils.Paths.Exists(instance.fileSurnames); b {
		_ = qb_utils.IO.ScanTextFromFile(instance.fileSurnames, func(_ string) bool {
			instance.maxSurnames++
			return false
		})
	}
	if b, _ := qb_utils.Paths.Exists(instance.fileCountryCodes); b {
		_ = qb_utils.IO.ScanTextFromFile(instance.fileCountryCodes, func(_ string) bool {
			instance.maxCountryCodes++
			return false
		})
	}
}

func (instance *GenUsersEngine) readName(i int) (name, gender string, err error) {
	count := 0
	var response string
	err = qb_utils.IO.ScanTextFromFile(instance.fileNames, func(text string) bool {
		count++
		if count == i {
			response = text
			return true
		}
		return false
	})
	if len(response) > 0 {
		tokens := strings.Split(strings.ToLower(response), ",")
		name = qb_utils.Strings.CapitalizeFirst(qb_utils.Arrays.GetAt(tokens, 0, "").(string))
		gender = qb_utils.Arrays.GetAt(tokens, 1, "").(string)
	}
	return
}

func (instance *GenUsersEngine) readSurname(i int) (surname string, err error) {
	count := 0
	var response string
	err = qb_utils.IO.ScanTextFromFile(instance.fileSurnames, func(text string) bool {
		count++
		if count == i {
			response = text
			return true
		}
		return false
	})
	if len(response) > 0 {
		surname = qb_utils.Strings.CapitalizeFirst(strings.ToLower(response))
	}
	return
}

func (instance *GenUsersEngine) readCountryCode(i int) (country, code, prefix string, err error) {
	count := 0
	var response string
	err = qb_utils.IO.ScanTextFromFile(instance.fileCountryCodes, func(text string) bool {
		count++
		if count == i {
			response = text
			return true
		}
		return false
	})
	if len(response) > 0 {
		tokens := qb_utils.Strings.Split(response, ",")
		country = strings.ReplaceAll(qb_utils.Arrays.GetAt(tokens, 0, "Italy").(string), "\"", "")
		code = qb_utils.Arrays.GetAt(tokens, 1, "IT").(string)
		prefix = qb_utils.Arrays.GetAt(tokens, 2, "+39").(string)
	}
	return
}
