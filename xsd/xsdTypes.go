package xsd

import "encoding/xml"

type Element struct {
	XMLName xml.Name `xml:"element"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type ComplexType struct {
	XMLName  xml.Name  `xml:"complexType"`
	Name     string    `xml:"name,attr"`
	Sequence []Element `xml:"sequence>element,omitempty"`
}

type Schema struct {
	XMLName      xml.Name      `xml:"schema"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
}
