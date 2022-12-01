package server

import (
	"github.com/achiku/xml"
	"io"
	"net/http"
	"regexp"
)

type ProcessARequest struct {
	XMLName   xml.Name `xml:"http://example.com/ns ProcessARequest"`
	RequestID string   `xml:"RequestID"`
}

type ProcessBRequest struct {
	XMLName   xml.Name `xml:"http://example.com/ns ProcessBRequest"`
	RequestID string   `xml:"RequestID"`
}

type ProcessAResponse struct {
	*AbstractResponse
	XMLName xml.Name `xml:"http://example.com/ns ProcessAResponse"`
	ID      string   `xml:"Id,omitifempty"`
	Process string   `xml:"Process,omitifempty"`
}

type ProcessBResponse struct {
	*AbstractResponse
	XMLName xml.Name `xml:"http://example.com/ns ProcessBResponse"`
	ID      string   `xml:"Id,omitifempty"`
	Process string   `xml:"Process,omitifempty"`
	Amount  string   `xml:"Amount,omitifempty"`
}

func processA() ProcessAResponse {
	return ProcessAResponse{
		AbstractResponse: &AbstractResponse{
			Code:   "200",
			Detail: "success",
		},
		ID:      "100",
		Process: "ProcessAResponse",
	}
}
func processB() ProcessBResponse {
	return ProcessBResponse{
		AbstractResponse: &AbstractResponse{
			Code:   "200",
			Detail: "success",
		},
		ID:      "100",
		Process: "ProcessBResponse",
		Amount:  "10000",
	}
}

func soapActionHandler(w http.ResponseWriter, r *http.Request) {
	soapAction := r.Header.Get("soapAction")
	var res interface{}
	switch soapAction {
	case "processA":
		res = processA()
	case "processB":
		res = processB()
	default:
		res = nil
	}
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: res,
		},
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml")
	x, err := xml.MarshalIndent(v, "", " ")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(x)
	return
}

func SoapBodyHandler(w http.ResponseWriter, r *http.Request) {
	rawbody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	a := regexp.MustCompile(`<ProcessARequest xmlns="http://example.com/ns>`)
	b := regexp.MustCompile(`<ProcessBRequest xmlns="http://example.com/ns>`)

	var res interface{}

	if a.MatchString(string(rawbody)) {
		res = processA()
	} else if b.MatchString(string(rawbody)) {
		res = processB()
	} else {
		res = nil
	}
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: res,
		},
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml")
	x, err := xml.MarshalIndent(v, "", " ")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(x)
	return
}

func NewSoapMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("dispatch/soapaction", soapActionHandler)
	mux.HandleFunc("dispatch/soapbody", soapBodyHandler)
	return mux
}

func NewSOAPServer(port string) *http.Server {
	mux := NewSOAPMux()
	server := &http.Server{
		Addr:    "localhost" + port,
		Handler: mux,
	}
	return server
}
