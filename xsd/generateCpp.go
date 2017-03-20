package xsd

type CppFile struct {
	FileName string
	Content  string
}

type generator interface {
	generateHeader() *CppFile
	generateSource() *CppFile
}

func (r *Restriction) ToCpp(typeName string) (header *CppFile, source *CppFile) {
	if r.Enumerations == nil {
		var generatorTpl generator

		switch r.Base {
		case "xs:string":
			generatorTpl = simpleTypeGenerator(typeName, "std::string")
		case "xs:positiveInteger":
			generatorTpl = simpleTypeGenerator(typeName, "unsigned int")
		case "xs:decimal":
			generatorTpl = simpleTypeGenerator(typeName, "double")
		case "xs:integer", "xs:int":
			generatorTpl = simpleTypeGenerator(typeName, "int")
		default:
		}

		if generatorTpl != nil {
			header = generatorTpl.generateHeader()
			source = generatorTpl.generateSource()
		}

		return
	}

	// in this case it's an enum

	enumTpl := enumGenerator(typeName, r.Enumerations)

	header = enumTpl.generateHeader()
	source = enumTpl.generateSource()

	return
}

func (st *SimpleType) ToCpp() (header *CppFile, source *CppFile) {
	return st.Restriction.ToCpp(st.Name)
}

func (ct *ComplexType) ToCpp() (header *CppFile, source *CppFile) {
	return
}

func (el *Element) ToCpp() (header *CppFile, source *CppFile) {
	return
}
