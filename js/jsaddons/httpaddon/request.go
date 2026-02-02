package httpaddon

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpRequest struct {
	Method string
	URL    string

	Header map[string][]string

	Body []byte
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{}
}

func (r *HttpRequest) SetMethod(method string) *HttpRequest {
	r.Method = method
	return r
}

func (r *HttpRequest) SetURL(url string) *HttpRequest {
	r.URL = url
	return r
}

func (r *HttpRequest) SetBodyString(body string) *HttpRequest {
	r.Body = []byte(body)
	return r
}

func (r *HttpRequest) SetBody(body []byte) *HttpRequest {
	r.Body = body
	return r
}
func (r *HttpRequest) SetBodyObj(body any) (*HttpRequest, error) {
	if r.Body == nil {
		return r, nil
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r.Body = bytes
	return r, nil
}

func (r *HttpRequest) SetBodyFormData(data FormDataCloser) *HttpRequest {
	header, body := data.close()
	return r.
		SetHeader("Content-Type", header).
		SetBody(body)
}

func (r *HttpRequest) SetHeader(key string, values ...string) *HttpRequest {
	if r.Header == nil {
		r.Header = make(map[string][]string)
	}
	r.Header[key] = values
	return r
}

func (r *HttpRequest) AddHeader(key string, values ...string) *HttpRequest {
	if r.Header == nil {
		r.Header = make(map[string][]string)
	}
	headers, ok := r.Header[key]
	if !ok {
		headers = values
	} else {
		headers = AddUniqueValues(headers, values...)
	}
	r.Header[key] = headers
	return r
}

func (r *HttpRequest) RemoveHeader(key string) *HttpRequest {
	delete(r.Header, key)
	return r
}

func (r *HttpRequest) Fetch() (*HttpResponse, error) {
	req, err := r.getGoRequest()
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	return ToHttpResponse(resp), err
}

func (r *HttpRequest) Get() (*HttpResponse, error) {
	r.Method = "GET"
	return r.Fetch()
}

func (r *HttpRequest) Post() (*HttpResponse, error) {
	r.Method = "POST"
	return r.Fetch()
}

func (r *HttpRequest) Put() (*HttpResponse, error) {
	r.Method = "PUT"
	return r.Fetch()
}

func (r *HttpRequest) Delete() (*HttpResponse, error) {
	r.Method = "DELETE"
	return r.Fetch()
}

func (r *HttpRequest) Head() (*HttpResponse, error) {
	r.Method = "HEAD"
	return r.Fetch()
}

func (r *HttpRequest) Options() (*HttpResponse, error) {
	r.Method = "OPTIONS"
	return r.Fetch()
}

func (r *HttpRequest) Connect() (*HttpResponse, error) {
	r.Method = "CONNECT"
	return r.Fetch()
}

func (r *HttpRequest) Trace() (*HttpResponse, error) {
	r.Method = "TRACE"
	return r.Fetch()
}

func (r *HttpRequest) Patch() (*HttpResponse, error) {
	r.Method = "PATCH"
	return r.Fetch()
}

func (r *HttpRequest) getGoRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	req.Header = r.Header

	return req, nil
}
