package httpaddon

import (
	"bytes"
	"mime/multipart"
)

type FormDataCloser interface {
	close() (headder string, body []byte)
}

type FormData struct {
	buffer *bytes.Buffer
	writer *multipart.Writer
}

func NewFormData() *FormData {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	return &FormData{
		buffer: buffer,
		writer: writer,
	}
}

func (f *FormData) AddField(key string, value string) *FormData {
	f.writer.WriteField(key, value)
	return f
}

func (f *FormData) AddFile(key string, filename string, body []byte) (*FormData, error) {
	part, err := f.writer.CreateFormFile(key, filename)
	if err != nil {
		return f, err
	}
	_, err = part.Write(body)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (f *FormData) close() (headder string, body []byte) {
	f.writer.Close()
	return f.writer.FormDataContentType(), f.buffer.Bytes()
}
