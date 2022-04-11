package pepper

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kz91/pepper/internal/utils"
)

type FuncMap map[string]interface{}

// 响应接口
type Response struct {
	Resp       http.ResponseWriter
	ErrorPages *ErrorPages
}

// 发送json格式数据
func (r *Response) Json(v interface{}) (err error) {
	r.SetHeader("Content-Type", "application/json")
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	_, err = r.Write(data)
	return
}

// 不返回错误信息的Json函数
func (r *Response) MustJson(code int, v interface{}) {
	if err := r.Json(v); err != nil {
		r.SetStatusCode(500)
	}
}

// 发送二进制内容
func (r *Response) Write(b []byte) (int, error) {
	return r.Resp.Write(b)
}

// 发送字符串
func (r *Response) WriteString(str string) (int, error) {
	return r.Write([]byte(str))
}

// 设置状态码
func (r *Response) SetStatusCode(code int) {
	r.Resp.WriteHeader(code)
}

// 设置响应头
func (r *Response) SetHeader(key, value string) {
	header := r.Resp.Header()
	header.Set(key, value)
}

// 设置cookie
func (r *Response) SetCookie(opt *http.Cookie) {
	http.SetCookie(r.Resp, opt)
}

// 发送文件
func (r *Response) WriteFile(file string, bufferSize int) (err error) {
	buffer := make([]byte, bufferSize)

	ext := filepath.Ext(file)
	mimeType := mime.TypeByExtension(ext)

	if mimeType == "" {
		mimeType = "text/plain"
	}

	r.SetHeader("Content-Type", mimeType+"; charset=utf-8")

	fp, err := os.Open(file)
	if err != nil {
		return
	}

	defer fp.Close()

	for {
		n := bufferSize
		n, err := fp.Read(buffer[:n])
		if err != nil && err != io.EOF {
			r.SetStatusCode(500)
			return err
		}

		r.Write(buffer[:n])
		if err == io.EOF {
			break
		}
	}

	return nil
}

// 发送错误页面并返回错误码
func (r *Response) SendErrorPage(code int) error {
	return (*ErrorPages).SendPage(r.ErrorPages, code, *r)
}

// 重定向
func (r *Response) Redirect(url string) {
	r.SetHeader("Location", url)
	r.SetStatusCode(302)
}

// 封装的 "html/template"
func (r *Response) Template(file string, tpl interface{}, fc FuncMap) error {
	t := template.New(utils.GetRandString(10))
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	t.Funcs(template.FuncMap(fc))

	t, err = t.Parse(string(data))
	if err != nil {
		return err
	}

	if err := t.Execute(r, tpl); err != nil {
		return err
	}

	return nil
}
