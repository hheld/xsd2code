package xsd

import "errors"

func (schema *Schema) FindType(typeName string) (tp TypePtr, err error) {
	for _, ct := range schema.ComplexTypes {
		if ct.Name == typeName {
			tp.Ct = &ct
			return
		}
	}

	for _, st := range schema.SimpleTypes {
		if st.Name == typeName {
			tp.St = &st
			return
		}
	}

	err = errors.New("Unknown type '" + typeName + "'")

	return
}
