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
		//fmt.Printf("%#v, %#v\n", ct.Name, ct.Attributes)
		ctf, _ := schema.FindType(ct.Name)
		fmt.Printf("complex type: %#v\n", ctf)
	}

	for _, st := range schema.SimpleTypes {
		//fmt.Printf("%#v\n", st)
		stf, _ := schema.FindType(st.Name)
		fmt.Printf("simple type: %#v\n", stf)
		st.ToCpp()
	}
}
