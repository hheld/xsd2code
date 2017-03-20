package xsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type simpleTypeTemplate struct {
	header string
	source string
	funcs  template.FuncMap
	values simpleTypeTemplateArgs
}

type simpleTypeTemplateArgs struct {
	TypeName       string
	SimpleTypeName string
}

func (stt *simpleTypeTemplate) generateHeader() *CppFile {
	headerTemplate, err := template.New("generateCppSimpleTypeHeader").Funcs(stt.funcs).Parse(stt.header)

	if err != nil {
		panic(err)
	}

	var headerFileContent bytes.Buffer

	err = headerTemplate.Execute(&headerFileContent, stt.values)

	if err != nil {
		panic(err)
	}

	headerFile := CppFile{
		FileName: fmt.Sprintf("%s.h", capitalizeFirst(stt.values.TypeName)),
		Content:  headerFileContent.String(),
	}

	return &headerFile
}

func (stt *simpleTypeTemplate) generateSource() *CppFile {
	sourceTemplate, err := template.New("generateCppSimpleTypeSource").Funcs(stt.funcs).Parse(stt.source)

	if err != nil {
		panic(err)
	}

	var sourceFileContent bytes.Buffer

	err = sourceTemplate.Execute(&sourceFileContent, stt.values)

	if err != nil {
		panic(err)
	}

	sourceFile := CppFile{
		FileName: fmt.Sprintf("%s.cpp", capitalizeFirst(stt.values.TypeName)),
		Content:  sourceFileContent.String(),
	}

	return &sourceFile
}

func simpleTypeGenerator(name, typeName string) generator {
	simpleTypeInstance := simpleTypeTemplate{
		header: simpleTypeHeaderTemplate,
		source: simpleTypeSourceTemplate,
		funcs: template.FuncMap{
			"toUpper":         strings.ToUpper,
			"capitalizeFirst": capitalizeFirst,
		},
		values: simpleTypeTemplateArgs{
			TypeName:       name,
			SimpleTypeName: typeName,
		},
	}

	return &simpleTypeInstance
}

const simpleTypeHeaderTemplate = `{{$includeGuardStr := .TypeName | toUpper | printf "%s_H" -}}
#ifndef {{$includeGuardStr}}
#define {{$includeGuardStr}}

{{if eq .SimpleTypeName "std::string" -}}
#include <string>

{{end -}}

{{$className := .TypeName | capitalizeFirst -}}
class {{$className}} final
{
public:
	{{$className}}();

	void setValue(const {{.SimpleTypeName}} &v);
	{{.SimpleTypeName}} value() const;

private:
	{{.SimpleTypeName}} value_;
};

#endif // {{$includeGuardStr}}
`

const simpleTypeSourceTemplate = `{{$className := .TypeName | capitalizeFirst -}}
#include "{{$className}}.h"

void {{$className}}::setValue(const {{.SimpleTypeName}} &v)
{
	value_ = v;
}

{{.SimpleTypeName}} {{$className}}::value() const
{
	return value_;
}
`
