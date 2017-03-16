package xsd

import (
	"encoding/xml"
	"math"
	"strconv"
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

func (st *SimpleType) isXsdType() {}

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Use     string   `xml:"use,attr,omitempty"`
}

type Element struct {
	XMLName   xml.Name `xml:"element"`
	Name      string   `xml:"name,attr"`
	Type      string   `xml:"type,attr"`
	MinOccurs int      `xml:"minOccurs,attr"`
	MaxOccurs int      `xml:"maxOccurs,attr"`
}

func (el *Element) isXsdType() {}

func (el *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// prevent recursion
	type Elem Element

	type element struct {
		Elem
		MaxOccursStr string `xml:"maxOccurs,attr"`
	}

	item := element{
		MaxOccursStr: "1",
		Elem: Elem{
			MinOccurs: 1,
		},
	}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	if item.MaxOccursStr == "unbounded" {
		el.MaxOccurs = int(math.MaxInt64)
	} else {
		el.MaxOccurs, _ = strconv.Atoi(item.MaxOccursStr)
	}

	el.Name = item.Name
	el.XMLName = item.XMLName
	el.Type = item.Type
	el.MinOccurs = item.MinOccurs

	return nil
}

type ComplexType struct {
	XMLName    xml.Name    `xml:"complexType"`
	Name       string      `xml:"name,attr"`
	Sequence   []Element   `xml:"sequence>element,omitempty"`
	Attributes []Attribute `xml:"attribute"`
}

func (ct *ComplexType) isXsdType() {}

type Schema struct {
	XMLName      xml.Name      `xml:"schema"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
}

type XsdType interface {
	isXsdType()
}
