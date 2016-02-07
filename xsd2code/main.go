package main

import (
	"encoding/xml"
	"fmt"
	"github.com/hheld/xsd2code/xsd"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("exampleSchema.xsd")

	if err != nil {
		panic(err)
	}

	schema := xsd.Schema{}

	err = xml.Unmarshal(b, &schema)

	if err != nil {
		panic(err)
	}

	for _, ct := range schema.ComplexTypes {
		ctf, _ := schema.FindType(ct.Name)
		fmt.Printf("complex type: %#v\n", ctf.Ct)
	}

	for _, st := range schema.SimpleTypes {
		stAsCpp := st.ToCpp()

		if stAsCpp.HeaderFile != nil {
			fmt.Printf("Header file %s:\n%s\n", stAsCpp.HeaderFile.FileName, stAsCpp.HeaderFile.Content)
		}

		if stAsCpp.SourceFile != nil {
			fmt.Printf("Source file %s:\n%s\n", stAsCpp.SourceFile.FileName, stAsCpp.SourceFile.Content)
		}

		if stAsCpp.SourceLine != nil {
			fmt.Printf("Source line:\n%s\n", *stAsCpp.SourceLine)
		}
	}
}
