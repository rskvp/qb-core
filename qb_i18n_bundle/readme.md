# G&G i18n #

![](icon.png)

G&G i18n is a Go package that helps you translate Go programs into multiple languages
using a [Mustache](https://mustache.github.io/) like syntax.

## Example ##

`en.json`: 

This is an i18n resource file.
This file contains some labels with both simple text and plurals.

NOTE: `{{rows-length}}` is not a model tag. 
This tag has been injected from LyGo Ext i18n.

```json
{
  "i18n-title": "A sample template",
  "i18n-description": "{{rows-index}}) Hello, my name is {{name}} {{surname}}",
  "i18n-rows": {
    "one": "There is {{rows-length}} row",
    "other": "There are {{rows-length}} rows"
  }
}
```

`model.json`: 

This is a sample data model passed to localizer.
The model contains some tags that are references to i18n resource file ("`{{i18n-description}}`").
Those tags will be localized and translated.  

```json
{
  "not_localized": "<h1>ROWS</h1>",
  "rows": [
    {
      "name": "Mario",
      "surname": "Rossi",
      "description": "{{i18n-description}}"
    },
    {
      "name": "Sergio",
      "surname": "Bianchi",
      "description": "{{i18n-description}}"
    }
  ]
}
```

`template.html`: 

This is a typical Mustache Template.
The template contains both i18 tags ("`{{i18n-title}}`, `{{i18n-rows}}`") and 
model tags ("`{{{not_localized}}}`", "`{{description}}`").
All this tags will be translated.

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{i18n-title}}</title>
</head>
<body>
    {{{not_localized}}}
    <h2>{{i18n-rows}}</h2>
    {{#rows}}
    <div>{{description}}</div>
    {{/rows}}
</body>
</html>
```

main.go:

```
package main

import (
	"bitbucket.org/lygo/lygo_commons/lygo_io"
	"bitbucket.org/lygo/lygo_commons/lygo_json"
	"bitbucket.org/lygo/lygo_ext_i18n/lygo_i18n"
	"fmt"
)

func main(){
    html, err := lygo_io.ReadTextFromFile("./template/template.html")
    if nil!=err{
        panic(err)
    }
    model, err := lygo_json.ReadMapFromFile("./template/model.json")
    if nil!=err{
    	panic(err)
    }

    // creates a i18n bundle with "it" as default language
    bundle, err := lygo_i18n.NewBundleFromDir("it", "./template")
    if nil!=err{
    	panic(err)
    }
    
    // creates a localizer to translate a text template
    localizer := lygo_i18n.NewLocalizer(bundle)
    localized, err := localizer.Localize("it", html, model)
    if nil != err {
        panic(err)
    }
    fmt.Println(localized)
}
```

output:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>A sample template</title>
</head>
<body>
    <h1>ROWS</h1>
    <h2>There are 2 rows</h2>
    <div>0) Hello, my name is Mario Rossi</div>
    <div>1) Hello, my name is Sergio Bianchi</div>
</body>
</html>
```
