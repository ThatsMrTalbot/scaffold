package encoding

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
)

// Responder is a response encoder
type Responder interface {
	Respond(status int, w http.ResponseWriter, response interface{}) error
}

// Parser is a request parser
type Parser interface {
	Parse(result interface{}, r *http.Request) error
}

// Encoding is a encoder
type Encoding interface {
	Mime() string
	Responder
	Parser
}

func readAll(r *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return buf, nil
}

type jsonResponder struct{}

func (*jsonResponder) Mime() string {
	return "application/json"
}

func (*jsonResponder) Respond(status int, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	buf, err := ffjson.Marshal(response)
	if err != nil {
		buf, _ = ffjson.Marshal(errorObj{Error: err.Error()})
	}

	w.WriteHeader(status)
	w.Write(buf)
	return err
}

func (*jsonResponder) Parse(result interface{}, r *http.Request) error {
	buf, err := readAll(r)
	if err != nil {
		return err
	}

	return ffjson.Unmarshal(buf, result)
}

type xmlResponder struct{}

func (*xmlResponder) Mime() string {
	return "application/xml"
}

func (*xmlResponder) Respond(status int, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/xml")

	buf, err := xml.Marshal(response)
	if err != nil {
		buf, _ = xml.Marshal(errorObj{Error: err.Error()})
	}

	w.WriteHeader(status)
	w.Write(buf)
	return err
}

func (*xmlResponder) Parse(result interface{}, r *http.Request) error {
	buf, err := readAll(r)
	if err != nil {
		return err
	}

	return xml.Unmarshal(buf, result)
}

// Responders & Responders
var (
	JSONEncoding Encoding = &jsonResponder{}
	XMLEncoding  Encoding = &xmlResponder{}
)
