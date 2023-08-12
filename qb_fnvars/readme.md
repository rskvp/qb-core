# Function Variables #
Expression engine for variables to solve at runtime.

This is useful when you need a dynamic value in template engines or text expressions.

Example: "`This is a RND number with 3 digits: <var>rnd|number|3</var>`"

```
func TestSolveAll(t *testing.T) {
	statements := []string{
		"this is a random string: (<var>rnd:alphanumeric:6</var>)",
		"this is a random number: (<var>rnd:numeric:6</var>)",
		"this is a random string LOWERCASE: (<var>rnd|alphanumeric|6|lower</var>)",
		"this is a random string UPPERCASE: (<var>rnd|alphanumeric|6|upper</var>)",
		"this is a random GUID UPPERCASE: (<var>rnd|guid|upper</var>)",
		"this is a random Num Between: (<var>rnd|between|0-10</var>)",
		"this is a random ID: (<var>rnd|id</var>)",
		"this is a DateTime: (<var>date|yyyy-MM-dd HH:mm|upper</var>)",
		"this is a DateTime ISO86001: (<var>date|iso|upper</var>)",
		"this is a DateTime Unix: (<var>date|unix</var>)",
		"this is a DateTime Ruby: (<var>date|ruby</var>)",
	}
	engine := gg.FnVars.NewEngine()
	for i, statement := range statements {
		text, err := engine.SolveText(statement)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println(i, statement, text)
	}

}
```

## Expressions Syntax ##

Function Variables are cosed in `<var>` tags.

Example: `<var>rnd</var>`

Each Function Variable can have some "options". Options are "optional" data
that allow a Function Variable to handle parametrization.

For example, if we need a 3 digits number we can write a Function Variable like this:
`<var>rnd|numeric|3</var>`. 

Here is Function Variable explanation:
 - `rnd`: Function Variable name
 - `numeric`: Option telling the engine to extract a number
 - `3`: Option telling the engine that the number should have 3 digits

