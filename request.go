package pepper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// 请求信息
type Request struct {
	Req           *http.Request
	Body          io.ReadCloser
	Method        string
	Path          string
	TrimPath      string
	Cookie        string
	Proto         string
	Host          string
	RemoteAddr    string
	ContentLength int64
	res           Response
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

func (r *Request) Scan(i interface{}) (err error) {
	var params map[string][]string

	iTypeOf := reflect.TypeOf(i)
	if iTypeOf.Kind() != reflect.Ptr {
		err = errors.New("parameter type must be point")
		return
	}

	iTypeOf = iTypeOf.Elem()
	iValueOf := reflect.ValueOf(i)

	if iTypeOf.Kind() != reflect.Struct {
		return
	}

	if r.Req.Method == "POST" {
		params = r.GetForm()
	} else {
		params = r.Req.URL.Query()
	}

	for pk, pv := range params {
		if len(pv) == 0 {
			continue
		}

		for i := 0; i < iTypeOf.NumField(); i++ {
			typeField := iTypeOf.Field(i)
			valueField := iValueOf.Elem().Field(i)
			name := typeField.Tag.Get("name")
			if name != pk {
				continue
			}

			if !valueField.CanSet() {
				err = errors.New("not settable: " + name)
				return
			}

			switch typeField.Type.Kind() {
			case reflect.Int64:
				var res int64
				res, err = strconv.ParseInt(pv[0], 10, 64)
				if err != nil {
					return
				}
				valueField.SetInt(res)

			case reflect.Float64:
				var res float64
				res, err = strconv.ParseFloat(pv[0], 64)
				if err != nil {
					return
				}
				valueField.SetFloat(res)

			case reflect.Bool:
				if pv[0] == "true" || pv[0] == "1" {
					valueField.SetBool(true)
					break
				}
				valueField.SetBool(false)

			case reflect.String:
				valueField.SetString(pv[0])
			}
		}

	}

	return
}
