package xsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type CppFileType struct {
	FileName string
	Content  string
}

type CppCodeType struct {
	SourceLine *string
	HeaderFile *CppFileType
	SourceFile *CppFileType
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
		var line string

		switch r.Base {
		case "xs:string":
			line = fmt.Sprintf("std::string %s_;", typeName)
		case "xs:positiveInteger":
			line = fmt.Sprintf("unsigned int %s_;", typeName)
		case "xs:decimal":
			line = fmt.Sprintf("double %s_;", typeName)
		case "xs:integer", "xs:int":
			line = fmt.Sprintf("int %s_;", typeName)
		default:
            line = fmt.Sprintf("std::string %s_;", typeName)
		}

		gen.SourceLine = &line
		return
	}

	// in this case it's an enum

	enumHeaderTemplate := `{{$includeGuardStr := .TypeName | toUpper | printf "%s_H"}}
#ifndef {{$includeGuardStr}}
#define {{$includeGuardStr}}

#include <string>
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

	enumSourceTemplate := `{{$enumName := .TypeName | capitalizeFirst}}
#include "{{$enumName}}.h"

namespace {{$enumName}}Conv
{
std::string toString({{$enumName}} v)
{
    std::string vAsStr;

    switch(v)
    {
    {{range .EnumValues}}
    case {{.Value}}:
        vAsStr = "{{.Value}}";
        break;
    {{end}}
    }

    return std::move(vAsStr);
}

{{$enumName}} fromString(const std::string &s)
{
    {{range .EnumValues}}
    if(s=="{{.Value}}") return {{.Value}};
    {{end}}
    throw("Unknown value '" + s + "' for enum '{{$enumName}}'.");
}
}

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

	headerTemplate, err := template.New("generateCppEnumHeader").Funcs(funcMap).Parse(enumHeaderTemplate)

	if err != nil {
		panic(err)
	}

	var headerFileContent bytes.Buffer

	templateValues := struct {
		TypeName   string
		EnumValues []Enumeration
	}{
		TypeName:   typeName,
		EnumValues: r.Enumerations,
	}

	err = headerTemplate.Execute(&headerFileContent, templateValues)

	if err != nil {
		panic(err)
	}

	headerFile := CppFileType{
		FileName: fmt.Sprintf("%s.h", capitalizeFirst(typeName)),
		Content:  headerFileContent.String(),
	}

	gen.HeaderFile = &headerFile

	sourceTemplate, err := template.New("generateCppEnumSource").Funcs(funcMap).Parse(enumSourceTemplate)

	if err != nil {
		panic(err)
	}

	var sourceFileContent bytes.Buffer

	err = sourceTemplate.Execute(&sourceFileContent, templateValues)

	if err != nil {
		panic(err)
	}

	sourceFile := CppFileType{
		FileName: fmt.Sprintf("%s.cpp", capitalizeFirst(typeName)),
		Content:  sourceFileContent.String(),
	}

	gen.SourceFile = &sourceFile

	return
}

func (st *SimpleType) ToCpp() CppCodeType {
	return st.Restriction.ToCpp(st.Name)
}
