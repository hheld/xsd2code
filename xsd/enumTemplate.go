package xsd

import (
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

func generateEnumTemplate(name string, values []Enumeration) enumTemplate {
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

	return enumInstance
}

const enumHeaderTemplate = `{{$includeGuardStr := .TypeName | toUpper | printf "%s_H" -}}
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

const enumSourceTemplate = `{{$enumName := .TypeName | capitalizeFirst -}}
#include "{{$enumName}}.h"

namespace {{$enumName}}Conv
{
std::string toString({{$enumName}} v)
{
    std::string vAsStr;

    switch(v)
    {
    {{- range .EnumValues}}
    case {{.Value}}:
        vAsStr = "{{.Value}}";
        break;
    {{- end}}
    }

    return std::move(vAsStr);
}

{{$enumName}} fromString(const std::string &s)
{
    {{- range .EnumValues}}
    if(s=="{{.Value}}") return {{.Value}};
    {{- end}}
    throw("Unknown value '" + s + "' for enum '{{$enumName}}'.");
}
}
`
