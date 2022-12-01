package client

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SoapHeader struct {
	XMLName xml.Name `xml:"x:Header"`
}

type SoapBody struct {
	XMLName xml.Name `xml:"x:Body"`
	Request interface{}
}

type SoapRoot struct {
	XMLName xml.Name `xml:"x:Envelope"`
	X       string   `xml:"x:xmlns:x,attr"`
	Sch     string   `xml:"x:xmlns:sch,attr"`
	Header  SoapHeader
	Body    SoapBody
}

type GetCitiesRequest struct {
	XMLName xml.Name `xml:"ns3.GetCitiesResponse"`
	result  struct{} `xml:result`
	cities  struct{} `xml:cities`
}

func SoapCall(service string, request interface{}) string {
	var root = SoapRoot{
		X:      "http://schemas.xmlsoap.org/soap/envelope/",
		Sch:    "http://www.n11mcom/ws/schemas",
		Header: SoapHeader{},
		Body: SoapBody{
			Request: request,
		},
	}
	out, _ := xml.MarshalIndent(&root, "", " ")
	body := string(out)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Post(service, "text/xml", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	return s
}
