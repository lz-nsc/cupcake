package cupcake

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Response struct {
	writer     http.ResponseWriter
	statusCode int
	render     *template.Template
}

func NewResponse(w http.ResponseWriter, render *template.Template) *Response {
	return &Response{
		writer: w,
		render: render,
	}
}

func (resp *Response) Status(code int) *Response {
	resp.statusCode = code
	resp.writer.WriteHeader(code)
	return resp
}

func (resp *Response) SetHeader(key string, value string) *Response {
	resp.writer.Header().Set(key, value)
	return resp
}

func (resp *Response) SetHeaders(headers map[string]string) *Response {
	for key, value := range headers {
		resp.writer.Header().Set(key, value)
	}
	return resp
}

func (resp *Response) String(code int, format string, values ...interface{}) {
	resp.SetHeader(
		"Content-Type", "text/plain",
	).Status(
		code,
	).write([]byte(fmt.Sprintf(format, values...)))
}

func (resp *Response) JSON(code int, obj interface{}) {
	resp.SetHeader("Content-Type", "application/json").Status(code)
	encoder := json.NewEncoder(resp.writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (resp *Response) Data(code int, data []byte) {
	resp.Status(code).write(data)
}

func (resp *Response) HTML(code int, html string) {
	resp.SetHeader(
		"Content-Type", "text/html",
	).Status(
		code,
	).write(
		[]byte(html),
	)
}

func (resp *Response) Render(code int, tmplName string, data interface{}) {
	resp.SetHeader(
		"Content-Type", "text/html",
	).Status(
		code,
	)

	if err := resp.render.ExecuteTemplate(resp.writer, tmplName, data); err != nil {
		resp.Error(http.StatusInternalServerError, err.Error())
	}
}

func (resp *Response) write(content []byte) {
	resp.writer.Write(content)
}

func (resp *Response) Error(errCode int, errMsg string) {
	http.Error(resp.writer, errMsg, errCode)
}

func (resp *Response) StatusCode() int {
	return resp.statusCode
}
