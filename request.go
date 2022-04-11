package pepper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// 请求信息
type Request struct {
	Req      *http.Request
	Method   string
	Path     string
	TrimPath string
	Cookie   string
	Proto    string
	Host     string
}

func (r *Request) GetHeader(key string) string {
	return r.Req.Header.Get(key)
}

func (r *Request) GetCookieValue(key string) string {
	cookie, err := r.Req.Cookie(key)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (r *Request) GetCookie(key string) (cookie *http.Cookie, err error) {
	cookie, err = r.Req.Cookie(key)
	return
}

func (r *Request) Query(key string) string {
	return r.Req.URL.Query().Get(key)
}

func (r *Request) GetForm() map[string][]string {
	headerContentType := r.GetHeader("Content-Type")
	contentType := strings.Split(headerContentType, ";")[0]
	if contentType != "application/x-www-form-urlencoded" {
		return nil
	}

	form := make(map[string][]string)

	body := r.Req.Body
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}

	formData := string(b)
	keyValueArray := strings.Split(formData, "/")

	for _, item := range keyValueArray {
		keyValue := strings.Split(item, "=")
		key := keyValue[0]
		value := keyValue[1]
		form[key] = append(form[key], value)
	}

	return form
}

func (r *Request) GetFormValue(key string) ([]string, bool) {
	form := r.GetForm()
	value, ok := form[key]
	return value, ok
}

func (r *Request) GetFormStringValue(key string) string {
	form, ok := r.GetFormValue(key)
	if !ok {
		return ""
	}

	if len(form) < 1 {
		return ""
	}

	return form[0]
}

func (r *Request) GetJson(i interface{}) error {
	headerContentType := r.GetHeader("Content-Type")
	contentType := strings.Split(headerContentType, ";")[0]
	if contentType != "application/json" {
		return nil
	}

	body := r.Req.Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, i)
}