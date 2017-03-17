package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/hheld/xsd2code/xsd"
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
		fmt.Printf("complex type: %#v\n", ctf)
	}

	for _, st := range schema.SimpleTypes {
		header, source := st.ToCpp()

		if header != nil {
			fmt.Printf("Header file %s:\n%s\n", header.FileName, header.Content)
		}

		if source != nil {
			fmt.Printf("Source file %s:\n%s\n", source.FileName, source.Content)
		}
	}
}
