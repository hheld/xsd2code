package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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

	const outDir = "cpp"

	if _, err := os.Stat("outDir"); os.IsNotExist(err) {
		err = os.Mkdir(outDir, 0755)
	}

	if err != nil {
		panic("There was an error: " + err.Error())
	}

	for _, ct := range schema.ComplexTypes {
		ctf, _ := schema.FindType(ct.Name)
		fmt.Printf("complex type: %#v\n", ctf)
	}

	for _, st := range schema.SimpleTypes {
		header, source := st.ToCpp()

		if header != nil {
			err = ioutil.WriteFile(path.Join(outDir, header.FileName), []byte(header.Content), 0644)

			if err != nil {
				panic("There was an error: " + err.Error())
			}
		}

		if source != nil {
			err = ioutil.WriteFile(path.Join(outDir, source.FileName), []byte(source.Content), 0644)

			if err != nil {
				panic("There was an error: " + err.Error())
			}
		}
	}
}
