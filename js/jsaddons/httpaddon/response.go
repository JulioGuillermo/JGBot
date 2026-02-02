package httpaddon

import (
	"io"
	"net/http"
)

type HttpResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0

	Header  map[string][]string
	Trailer map[string][]string

	ContentLength int64
	body          io.ReadCloser

	TransferEncoding []string
	Close            bool
	Uncompressed     bool
}

func ToHttpResponse(response *http.Response) *HttpResponse {
	if response == nil {
		return nil
	}
	return &HttpResponse{
		Status:           response.Status,
		StatusCode:       response.StatusCode,
		Proto:            response.Proto,
		ProtoMajor:       response.ProtoMajor,
		ProtoMinor:       response.ProtoMinor,
		Header:           response.Header,
		Trailer:          response.Trailer,
		ContentLength:    response.ContentLength,
		body:             response.Body,
		TransferEncoding: response.TransferEncoding,
		Close:            response.Close,
		Uncompressed:     response.Uncompressed,
	}
}

func (r *HttpResponse) BodyBytes() ([]byte, error) {
	defer r.CloseBody()
	return io.ReadAll(r.body)
}

func (r *HttpResponse) BodyString() (string, error) {
	bytes, err := r.BodyBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *HttpResponse) CloseBody() {
	r.body.Close()
}
