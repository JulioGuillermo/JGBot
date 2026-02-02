package httpaddon

import (
	"JGBot/js/exec"

	"github.com/fastschema/qjs"
)

func WithHttp() exec.Option {
	return func(e *exec.Executor) error {
		e.AddFunc("HttpRequest", func(this *qjs.This) (*qjs.Value, error) {
			return qjs.ToJsValue(this.Context(), NewHttpRequest())
		})
		e.AddFunc("HttpFormData", func(this *qjs.This) (*qjs.Value, error) {
			return qjs.ToJsValue(this.Context(), NewFormData())
		})
		return nil
	}
}
