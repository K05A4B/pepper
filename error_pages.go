package pepper

import (
	"net/http"
	"strconv"
)

type ErrorPages struct {
	NotFound            string
	Forbidden           string
	InternalServerError string
	Other               map[int]string
}

func (e *ErrorPages) SendPage(code int, res Response) error {
	msg := http.StatusText(code)
	page := ""
	switch code {
	case 404:
		page = e.NotFound
	case 403:
		page = e.Forbidden
	case 500:
		page = e.InternalServerError
	default:
		var ok bool
		page, ok = e.Other[code]
		if !ok {
			page = ""
		}
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
	e.Other[code] = file
}