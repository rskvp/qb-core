package npm

import (
	"bytes"
	"strings"
	"text/template"
)

type DataPackage struct {
	Name            string
	Description     string
	Version         string
	Main            string
	Scripts         string
	Dependencies    string
	DevDependencies string
	RepositoryURL   string
	Author          string
	License         string
}

const tplPackage = `
{
  "name": "{{.Name}}",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "{{.Main}}",
  "scripts": {
    {{.Scripts}}
  },
  "dependencies": {
    {{.Dependencies}}
  },
  "devDependencies": {
    {{.DevDependencies}}
  },
  "repository": {
    "type": "git",
    "url": "{{.RepositoryURL}}"
  },
  "author": "{{.Author}}",
  "license": "{{.License}}"
}
`

func MergeTpl(text string, data *DataPackage) (string, error) {
	buff := bytes.NewBufferString("")
	t := template.Must(template.New("package").Parse(strings.Trim(text, "\n")))
	err := t.Execute(buff, data)
	if nil != err {
		return "", err
	}
	return buff.String(), err
}
