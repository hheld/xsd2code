package xsd

import "errors"

func (schema *Schema) FindType(typeName string) (tp XsdType, err error) {
	for _, ct := range schema.ComplexTypes {
		if ct.Name == typeName {
			tp = &ct
			return
		}
	}

	for _, st := range schema.SimpleTypes {
		if st.Name == typeName {
			tp = &st
			return
		}
	}

	for _, el := range schema.Elements {
		if el.Name == typeName {
			tp = &el
			return
		}
	}

	err = errors.New("Unknown type '" + typeName + "'")

	return
}
