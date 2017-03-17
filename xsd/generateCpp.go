package xsd

import (
	"bytes"
	"fmt"
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

	enumTpl := generateEnumTemplate(typeName, r.Enumerations)

	headerTemplate, err := template.New("generateCppEnumHeader").Funcs(enumTpl.funcs).Parse(enumTpl.header)

	if err != nil {
		panic(err)
	}

	var headerFileContent bytes.Buffer

	err = headerTemplate.Execute(&headerFileContent, enumTpl.values)

	if err != nil {
		panic(err)
	}

	headerFile := CppFileType{
		FileName: fmt.Sprintf("%s.h", capitalizeFirst(typeName)),
		Content:  headerFileContent.String(),
	}

	gen.HeaderFile = &headerFile

	sourceTemplate, err := template.New("generateCppEnumSource").Funcs(enumTpl.funcs).Parse(enumTpl.source)

	if err != nil {
		panic(err)
	}

	var sourceFileContent bytes.Buffer

	err = sourceTemplate.Execute(&sourceFileContent, enumTpl.values)

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
