package upload

import (
	"errors"
	"mime"
	"mime/multipart"

	"github.com/kz91/pepper"
)

// 创建一个接收器
func NewReceive(r *pepper.Request, key string, rule *Rule) (recvFiles *Files, err error) {
	if r.Method != "POST" {
		err = errors.New("method is not post")
		return
	}

	mediaType, params, err := mime.ParseMediaType(r.GetHeader("Content-Type"))
	if err != nil {
		return
	}

	if mediaType != "multipart/form-data" {
		err = errors.New("content-type is not \"multipart/form-data\"")
		return
	}

	boundary, existBoundary := params["boundary"]
	if !existBoundary {
		return
	}

	reader := multipart.NewReader(r.Req.Body, boundary)

	recvFiles = &Files{
		reader:  reader,
		rule:    rule,
		key:     key,
		partNum: 0,
	}

	return
}