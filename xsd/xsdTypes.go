package xsd

import (
	"encoding/xml"
)

type Enumeration struct {
	XMLName xml.Name `xml:"enumeration"`
	Value   string   `xml:"value,attr"`
}

type Pattern struct {
	XMLName xml.Name `xml:"pattern"`
	Value   string   `xml:"value,attr"`
}

type Restriction struct {
	XMLName      xml.Name      `xml:"restriction"`
	Base         string        `xml:"base,attr"`
	Enumerations []Enumeration `xml:"enumeration"`
	Patterns     Pattern       `xml:"pattern"`
}

type SimpleType struct {
	XMLName     xml.Name    `xml:"simpleType"`
	Name        string      `xml:"name,attr"`
	Restriction Restriction `xml:"restriction"`
}

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
	SimpleTypes  []SimpleType  `xml:"simpleType"`
}

type TypePtr struct {
	Ct *ComplexType
	St *SimpleType
}
