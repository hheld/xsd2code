package xsd

import "encoding/xml"

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Use     string   `xml:"use,attr,omitempty"`
}

type Element struct {
	XMLName xml.Name `xml:"element"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type ComplexType struct {
	XMLName    xml.Name    `xml:"complexType"`
	Name       string      `xml:"name,attr"`
	Sequence   []Element   `xml:"sequence>element,omitempty"`
	Attributes []Attribute `xml:"attribute"`
}

type Schema struct {
	XMLName      xml.Name      `xml:"schema"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
}
