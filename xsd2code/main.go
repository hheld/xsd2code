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

	fmt.Printf("%#v\n", schema)

	for _, ct := range schema.ComplexTypes {
		fmt.Printf("%#v\n", ct.Name)
	}
}
