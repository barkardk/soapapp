package server

import "github.com/achiku/xml"

type SOAPEnvelope struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Envelope"`
	Header  *SOAPHeader `xml:"omitempty"`
	Body    *SOAPBody   `xml:"omitempty"`
}

type SOAPHeader struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Header"`
	Content interface{} `xml:"omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Body"`
	Fault   *SOAPFault  `xml:"omitempty"`
	Content interface{} `xml:"omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope Body"`
	Code    string   `xml:"faultcode,omitempty"`
	String  string   `xml:"faultstring,omitempty"`
	Actor   string   `xml:"faultactor,omitempty"`
	Detail  string   `xml:"detail,omitempty"`
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}
	var (
		token xml.Token
		err	error
		consumed bool
	)
	Loop:
		for {
			if token, err = d.Token(); err != nil {
				return err
			}
			if token == nil {
				break
			}
			envelopeNameSpace := "http://schemas.xmlsoap.org/soap/envelope/"
			switch se := token.(type) {
			case xml.StartElement:
				if consumed {
					return xml.UnmarshalError("Found muliple elements inside SOAP body: not wrapped-document/literal WS-I compliant")
				}else if
			}
		}
	return nil
}

type Name struct {
}
type Auth struct {
}

type Person struct {
}
type AbstractResponse struct {
	Code   string `xml:"Code:omitempty"`
	Detail string `xml:"Detail:omitempty"`
}
type ConcreteResponse struct {
}
