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
		switch r.Base {
		case "xs:string":
		case "xs:positiveInteger":
		case "xs:decimal":
		case "xs:integer", "xs:int":
		default:
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
