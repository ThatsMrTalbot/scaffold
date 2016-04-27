package encoding

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/ThatsMrTalbot/scaffold"
	"golang.org/x/net/context"
)

// Defaults for ease of use
var (
	DefaultEncoder        = NewEncoder(JSONEncoding, XMLEncoding)
	DefaultHandlerBuilder = DefaultEncoder.HandlerBuilder
)

// Encoder encodes and decodes requests based on the mime type
type Encoder struct {
	defaultEncoding Encoding
	encodings       map[Encoding]*mimeInfo
}

// NewEncoder creates encoder
func NewEncoder(defaultEncoding Encoding, encodings ...Encoding) *Encoder {
	m := make(map[Encoding]*mimeInfo)

	for _, e := range append(encodings, defaultEncoding) {
		info := parseMimeInfo(e.Mime())
		m[e] = info
	}

	return &Encoder{
		defaultEncoding: defaultEncoding,
		encodings:       m,
	}
}

// Parser returns a parser based on the mime, if none can be matched the
// default is returned
func (e *Encoder) Parser(r *http.Request) Parser {
	contentType := r.Header.Get("Content-Type")
	info := parseMimeInfo(contentType)

	if info != nil {
		for encoding, i := range e.encodings {
			if i.equals(info) {
				return encoding
			}
		}
	}

	return e.defaultEncoding
}

// Responder returns a responder based on the mime, if none can be matched the
// default is returned
func (e *Encoder) Responder(r *http.Request) Responder {
	accept := r.Header.Get("Accept")
	infoList := parseAcceptHeader(accept)

	for _, info := range infoList {
		for encoding, i := range e.encodings {
			if i.equals(info) {
				return encoding
			}
		}
	}

	return e.Parser(r).(Responder)
}

// HandlerBuilder can be used in scaffold to create handlers based on function
func (s *Selector) HandlerBuilder(handler interface{}) (scaffold.Handler, error) {
	typ := reflect.TypeOf(handler)

	if typ.Kind() != reflect.Func {
		return nil, errors.New("Invalid handler, must be a function")
	}

	if typ.NumIn() == 0 {
		return nil, errors.New("Invalid handler, must have at lease one argument")
	}

	params := make([]func(context.Context, http.ResponseWriter, *http.Request) (reflect.Value, error), typ.NumIn())

	reqTyp := typ.In(0)
	if reqTyp.Kind() != reflect.Struct {
		return nil, errors.New("Invalid handler, first argument must be a struct")
	}

	params[0] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
		parser := s.Parser(r)
		val := reflect.New(reqTyp)
		err := parser.Parse(val.Interface(), r)
		return val.Elem(), err
	}

	for i := 1; i < typ.NumIn(); i++ {
		switch typ.In(i) {
		case reflect.TypeOf((*context.Context)(nil)).Elem():
			params[i] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(ctx), nil
			}
		case reflect.TypeOf((*http.ResponseWriter)(nil)).Elem():
			params[i] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(w), nil
			}
		case reflect.TypeOf((*http.Request)(nil)):
			params[i] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(r), nil
			}
		default:
			return nil, errors.New("Invalid handler, invalid argument type")
		}
	}

	if typ.NumIn() != 2 {
		return nil, errors.New("Invalid handler, must have at two return parameters")
	}

	if typ.Out(0).Kind() != reflect.Struct {
		return nil, errors.New("Invalid handler, first return parameter must be a struct")
	}

	if typ.Out(0) != reflect.TypeOf((error)(nil)) {
		return nil, errors.New("Invalid handler, second return parameter must be an error")
	}

	return &caller{
		s:      s,
		params: params,
		caller: reflect.ValueOf(handler),
	}, nil
}
