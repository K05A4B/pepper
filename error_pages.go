package pepper

import (
	"net/http"
	"strconv"
)

type ErrorPages struct {
	Static struct {
		NotFound            string
		Forbidden           string
		InternalServerError string
		Other               map[int]string
	}
	NotFound            HandlerFunc
	Forbidden           HandlerFunc
	InternalServerError HandlerFunc
	Other               map[int]HandlerFunc
}

func (e *ErrorPages) SendPage(code int, res Response, req *Request) error {
	var handler HandlerFunc
	msg := http.StatusText(code)
	page := ""

	switch code {
	case 404:
		page = e.Static.NotFound
		handler = e.NotFound
	case 403:
		page = e.Static.Forbidden
		handler = e.Forbidden
	case 500:
		page = e.Static.InternalServerError
		handler = e.InternalServerError
	default:
		var ok bool
		page, ok = e.Static.Other[code]
		if !ok {
			page = ""
		}

		handler = e.Other[code]
	}

	if handler != nil {
		res.SetStatusCode(code)
		handler(res, req)
		return nil
	}

	if msg == "" {
		return res.SendErrorPage(500)
	}

	res.SetStatusCode(code)

	if page == "" {
		res.WriteString(strconv.Itoa(code) + " " + msg)
		return nil
	}

	return res.WriteFile(page, 5120)

}

func (e *ErrorPages) Set(code int, file string) {
	e.Static.Other[code] = file
}

func (e *ErrorPages) SetHandler(code int, handler HandlerFunc) {
	e.Other[code] = handler
}

func NewErrorPages() *ErrorPages {
	e := &ErrorPages{
		Other: make(map[int]HandlerFunc),
	}

	e.Static.Other = make(map[int]string)
	return e
}
