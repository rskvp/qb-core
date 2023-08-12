package qb_utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type CsvHelper struct {
}

var CSV *CsvHelper

func init() {
	CSV = new(CsvHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type CsvOptions struct {
	Comma          string `json:"comma"`
	Comment        string `json:"comment"`
	FirstRowHeader bool   `json:"first_row_header"`
}

func (instance *CsvHelper) NewCsvOptions(comma string, comment string, firstRowHeader bool) *CsvOptions {
	return &CsvOptions{
		Comma:          comma,
		Comment:        comment,
		FirstRowHeader: firstRowHeader,
	}
}

func (instance *CsvHelper) NewCsvOptionsDefaults() *CsvOptions {
	return &CsvOptions{
		Comma:          ";",
		Comment:        "#",
		FirstRowHeader: true,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *CsvHelper) ReadAll(in string, options *CsvOptions) (response []map[string]string, err error) {

	records, err := readRecords(in, options)
	if err != nil {
		return nil, err
	}

	response = make([]map[string]string, 0)
	headers := buildHeaders(&records, options)
	for _, row := range records {
		item := make(map[string]string)
		for i, value := range row {
			if len(headers) > i {
				item[headers[i]] = value
			}
		}
		response = append(response, item)
	}

	return response, err
}

func (instance *CsvHelper) WriteAll(data []map[string]interface{}, options *CsvOptions) (response string, err error) {
	buf := bytes.NewBufferString("")
	w := csv.NewWriter(buf)
	w.Comma = []rune(options.Comma)[0]
	w.UseCRLF = false
	defer w.Flush()

	err = write(w, data, options.FirstRowHeader, Maps.Keys(data[0]))
	return buf.String(), err
}

func (instance *CsvHelper) WriteFile(data []map[string]interface{}, options *CsvOptions, filename string) (err error) {
	// remove file if exists
	_ = IO.Remove(filename)

	return instance.AppendFile(data, options, filename)
}

func (instance *CsvHelper) AppendFile(data []map[string]interface{}, options *CsvOptions, filename string) (err error) {
	size, _ := IO.FileSize(filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	w.Comma = []rune(options.Comma)[0]
	w.UseCRLF = false
	defer w.Flush()

	// handle header using fields order like already existing
	addHeader := options.FirstRowHeader && size == 0
	header := Maps.Keys(data[0]) // keys are sorted randomly
	if size > 0 && options.FirstRowHeader{
		// get header
		text := IO.ReadLinesFromFile(filename, 1)
		if len(text)>0{
			records, err := readRecords(text, options)
			if nil==err && len(records)>0{
				header = records[0] // existing sorted fields
			}
		}
	}
	return write(w, data, addHeader, header)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func setOptions(r *csv.Reader, options *CsvOptions) {
	if nil != r && nil != options {
		if len(options.Comma) == 1 {
			r.Comma = []rune(options.Comma)[0]
		}
		if len(options.Comment) == 1 {
			r.Comment = []rune(options.Comment)[0]
		}
	}
}

func buildHeaders(records *[][]string, options *CsvOptions) []string {
	headers := make([]string, 0)
	if options.FirstRowHeader && len(*records) > 1 {
		headers = (*records)[0]
		*records = (*records)[1:][:]
	} else {
		for i := 0; i < len(*records); i++ {
			headers = append(headers, fmt.Sprintf("field_%v", i))
		}
	}

	return headers
}

func write(w *csv.Writer, data []map[string]interface{}, firstRowHeader bool, header []string) error {
	rows := make([][]string, 0)

	// write header
	if firstRowHeader {
		rows = append(rows, header)
	}

	// write rows
	for _, m := range data {
		rows = append(rows, toArrayOfQuotedStrings(Maps.ValuesOfKeys(m, header)))
	}

	return w.WriteAll(rows)
}

func toArrayOfQuotedStrings(values []interface{}) []string {
	response := make([]string, 0)
	for _, value := range values {
		response = append(response, Convert.ToString(value))
	}
	return response
}

func readRecords(text string, options *CsvOptions)([][]string, error){
	r := csv.NewReader(strings.NewReader(text))
	setOptions(r, options)
	return r.ReadAll()
}