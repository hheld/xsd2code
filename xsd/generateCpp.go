package xsd

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func capitalizeFirst(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(string(s[0])) + s[1:]
	} else if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}

	// s = ""
	return s
}

func (r *Restriction) ToCpp(typeName string) {
	if r.Enumerations == nil {
		fmt.Printf("%s\n", r.Base)
		return
	}

	// in this case it's an enum

	enumTemplate := `
#ifndef {{.TypeName | toUpper}}_H
#define {{.TypeName | toUpper}}_H

enum class {{.TypeName | capitalizeFirst}}
{
    {{.EnumValues | enumToString}}
};

#endif // {{.TypeName | toUpper}}_H

`

	funcMap := template.FuncMap{
		"toUpper":         strings.ToUpper,
		"capitalizeFirst": capitalizeFirst,
		"enumToString": func(enumValues []Enumeration) string {
			s := make([]string, len(enumValues))

			for i, ev := range enumValues {
				s[i] = ev.Value
			}

			return strings.Join(s, ",\n    ")
		},
	}

	tmpl, err := template.New("generateCppEnum").Funcs(funcMap).Parse(enumTemplate)

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		TypeName   string
		EnumValues []Enumeration
	}{
		TypeName:   typeName,
		EnumValues: r.Enumerations,
	})

	if err != nil {
		panic(err)
	}
}

func (st *SimpleType) ToCpp() {
	st.Restriction.ToCpp(st.Name)
}
