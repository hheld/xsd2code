package xsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type complexTypeTemplate struct {
	header string
	source string
	funcs  template.FuncMap
	values complexTypeTemplateArgs
}

type complexTypeTemplateArgs struct {
	TypeName   string
	Elements   []Element
	Attributes []Attribute
}

func (ctt *complexTypeTemplate) generateHeader() *CppFile {
	headerTemplate, err := template.New("generateCppComplexTypeHeader").Funcs(ctt.funcs).Parse(ctt.header)

	if err != nil {
		panic(err)
	}

	var headerFileContent bytes.Buffer

	err = headerTemplate.Execute(&headerFileContent, ctt.values)

	if err != nil {
		panic(err)
	}

	headerFile := CppFile{
		FileName: fmt.Sprintf("%s.h", capitalizeFirst(ctt.values.TypeName)),
		Content:  headerFileContent.String(),
	}

	return &headerFile
}

func (ctt *complexTypeTemplate) generateSource() *CppFile {
	sourceTemplate, err := template.New("generateCppComplexTypeSource").Funcs(ctt.funcs).Parse(ctt.source)

	if err != nil {
		panic(err)
	}

	var sourceFileContent bytes.Buffer

	err = sourceTemplate.Execute(&sourceFileContent, ctt.values)

	if err != nil {
		panic(err)
	}

	sourceFile := CppFile{
		FileName: fmt.Sprintf("%s.cpp", capitalizeFirst(ctt.values.TypeName)),
		Content:  sourceFileContent.String(),
	}

	return &sourceFile
}

func complexTypeGenerator(name string, attributes []Attribute, elements []Element) generator {
	complexTypeInstance := complexTypeTemplate{
		header: complexTypeHeaderTemplate,
		source: complexTypeSourceTemplate,
		funcs: template.FuncMap{
			"toUpper":         strings.ToUpper,
			"capitalizeFirst": capitalizeFirst,
		},
		values: complexTypeTemplateArgs{
			TypeName:   name,
			Attributes: attributes,
			Elements:   elements,
		},
	}

	return &complexTypeInstance
}

const complexTypeHeaderTemplate = `{{$includeGuardStr := .TypeName | toUpper | printf "%s_H" -}}
#ifndef {{$includeGuardStr}}
#define {{$includeGuardStr}}

{{range .Attributes -}}
#include "{{.Type | capitalizeFirst}}.h"
{{end}}
{{range .Elements -}}
#include "{{.Type | capitalizeFirst}}.h"
{{end}}

{{$className := .TypeName | capitalizeFirst -}}
class {{$className}} final
{
public:
	{{$className}}();

	{{range .Attributes -}}
	const {{.Type | capitalizeFirst}} &{{.Name}}() const;
	void set{{.Name | capitalizeFirst}}(const {{.Type | capitalizeFirst}} &v);

	{{end}}
	{{range .Elements -}}
	const {{.Type | capitalizeFirst}} &{{.Name}}() const;
	void set{{.Name | capitalizeFirst}}(const {{.Type | capitalizeFirst}} &v);

	{{end}}
private:
	{{range .Attributes -}}
	{{.Type | capitalizeFirst}} {{.Name}}_;
	{{end}}
	{{range .Elements -}}
	{{.Type | capitalizeFirst}} {{.Name}}_;
	{{end}}
};

#endif // {{$includeGuardStr}}`

const complexTypeSourceTemplate = `{{$className := .TypeName | capitalizeFirst -}}
#include "{{$className}}.h"

{{$className}}::{{$className}}()
{}

{{range .Attributes -}}
const {{.Type | capitalizeFirst}} &{{$className}}::{{.Name}}() const
{
	return {{.Name}}_;
}

void {{$className}}::set{{.Name | capitalizeFirst}}(const {{.Type | capitalizeFirst}} &v)
{
	{{.Name}}_ = v;
}

{{end}}
{{range .Elements -}}
const {{.Type | capitalizeFirst}} &{{$className}}::{{.Name}}() const
{
	return {{.Name}}_;
}

void {{$className}}::set{{.Name | capitalizeFirst}}(const {{.Type | capitalizeFirst}} &v)
{
	{{.Name}}_ = v;
}

{{end}}
`
