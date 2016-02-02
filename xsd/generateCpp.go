package xsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type CppCodeType struct {
	SourceLine *string
	HeaderFile *string
	SourceFile *string
}

func capitalizeFirst(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(string(s[0])) + s[1:]
	} else if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}

	// s = ""
	return s
}

func (r *Restriction) ToCpp(typeName string) (gen CppCodeType) {
	if r.Enumerations == nil {
		fmt.Printf("%s\n", r.Base)
        switch r.Base {
        case "xs:string":
            line := fmt.Sprintf("std::string %s_;", typeName)
            gen.SourceLine = &line
            return
        }
		return
	}

	// in this case it's an enum

	enumTemplate := `
{{$includeGuardStr := .TypeName | toUpper | printf "%s_H"}}
#ifndef {{$includeGuardStr}}
#define {{$includeGuardStr}}
{{$enumName := .TypeName | capitalizeFirst}}
enum class {{$enumName}}
{
    {{.EnumValues | enumToString}}
};

namespace {{$enumName}}Conv
{
std::string toString({{$enumName}} v);
{{$enumName}} fromString(const std::string &s);
}

#endif // {{$includeGuardStr}}

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

	var headerFile bytes.Buffer

	err = tmpl.Execute(&headerFile, struct {
		TypeName   string
		EnumValues []Enumeration
	}{
		TypeName:   typeName,
		EnumValues: r.Enumerations,
	})

	if err != nil {
		panic(err)
	}

	headerFileStr := headerFile.String()
	gen.HeaderFile = &headerFileStr

	return
}

func (st *SimpleType) ToCpp() CppCodeType {
	return st.Restriction.ToCpp(st.Name)
}
