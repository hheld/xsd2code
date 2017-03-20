package xsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type enumTemplate struct {
	header string
	source string
	funcs  template.FuncMap
	values enumTemplateArgs
}

type enumTemplateArgs struct {
	TypeName   string
	EnumValues []Enumeration
}

func (et *enumTemplate) generateHeader() *CppFile {
	headerTemplate, err := template.New("generateCppEnumHeader").Funcs(et.funcs).Parse(et.header)

	if err != nil {
		panic(err)
	}

	var headerFileContent bytes.Buffer

	err = headerTemplate.Execute(&headerFileContent, et.values)

	if err != nil {
		panic(err)
	}

	headerFile := CppFile{
		FileName: fmt.Sprintf("%s.h", capitalizeFirst(et.values.TypeName)),
		Content:  headerFileContent.String(),
	}

	return &headerFile
}

func (et *enumTemplate) generateSource() *CppFile {
	sourceTemplate, err := template.New("generateCppEnumSource").Funcs(et.funcs).Parse(et.source)

	if err != nil {
		panic(err)
	}

	var sourceFileContent bytes.Buffer

	err = sourceTemplate.Execute(&sourceFileContent, et.values)

	if err != nil {
		panic(err)
	}

	sourceFile := CppFile{
		FileName: fmt.Sprintf("%s.cpp", capitalizeFirst(et.values.TypeName)),
		Content:  sourceFileContent.String(),
	}

	return &sourceFile
}

func enumGenerator(name string, values []Enumeration) generator {
	enumInstance := enumTemplate{
		header: enumHeaderTemplate,
		source: enumSourceTemplate,
		funcs: template.FuncMap{
			"toUpper":         strings.ToUpper,
			"capitalizeFirst": capitalizeFirst,
			"enumToString": func(enumValues []Enumeration) string {
				s := make([]string, len(enumValues))

				for i, ev := range enumValues {
					s[i] = ev.Value
				}

				return strings.Join(s, ",\n    ")
			},
		},
		values: enumTemplateArgs{
			TypeName:   name,
			EnumValues: values,
		},
	}

	return &enumInstance
}

const enumHeaderTemplate = `{{$includeGuardStr := .TypeName | toUpper | printf "%s_H" -}}
#ifndef {{$includeGuardStr}}
#define {{$includeGuardStr}}

#include <string>
{{$enumName := (print .TypeName "Enum") | capitalizeFirst}}
enum class {{$enumName}}
{
    {{.EnumValues | enumToString}}
};

namespace {{$enumName}}Conv
{
std::string toString({{$enumName}} v);
{{$enumName}} fromString(const std::string &s);
}

{{$className := .TypeName | capitalizeFirst -}}
class {{$className}} final
{
public:
	{{$className}}();

	void setValue(const std::string &v);
	std::string value() const;

private:
	{{$enumName}} value_;
};

#endif // {{$includeGuardStr}}
`

const enumSourceTemplate = `{{$enumName := (print .TypeName "Enum") | capitalizeFirst -}}
#include "{{.TypeName | capitalizeFirst}}.h"

namespace {{$enumName}}Conv
{
std::string toString({{$enumName}} v)
{
    std::string vAsStr;

    switch(v)
    {
    {{- range .EnumValues}}
    case {{$enumName}}::{{.Value}}:
        vAsStr = "{{.Value}}";
        break;
    {{- end}}
    }

    return std::move(vAsStr);
}

{{$enumName}} fromString(const std::string &s)
{
    {{- range .EnumValues}}
    if(s=="{{.Value}}") return {{$enumName}}::{{.Value}};
    {{- end}}
    throw("Unknown value '" + s + "' for enum '{{$enumName}}'.");
}
}

{{$className := .TypeName | capitalizeFirst -}}
{{$className}}::{{$className}}()
{

}

void {{$className}}::setValue(const std::string &v)
{
	value_ = {{$enumName}}Conv::fromString(v);
}

std::string {{$className}}::value() const
{
	return {{$enumName}}Conv::toString(value_);
}
`
