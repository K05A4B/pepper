package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
)

type Files struct {
	reader  *multipart.Reader
	rule    *Rule
	key     string
	partNum int
}

func (fs *Files) Next() (f *File, err error) {

	if fs == nil {
		err = errors.New("invalid memory address or nil pointer dereference")
		return
	}

	var part *multipart.Part
	rule := fs.rule

	for {
		part, err = fs.reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		name := part.FileName()
		form := part.FormName()
		mime := part.Header.Get("Content-Type")

		// 判断是不是这个form
		if form != fs.key {
			continue
		}

		if (fs.partNum + 1) > fs.rule.MaxNumber && fs.rule.MaxNumber != 0 {
			err = errors.New("maximum number exceeded")
			return
		}

		fs.partNum++

		// 检查 mime 是否允许 &&
		if !rule.Mime.Exist(mime) && !rule.Mime.Exist(filepath.Ext(name)) &&
			rule.Mime.Len() != 0 {
			err = errors.New("this file type is not supported")
			return
		}

		f = &File{
			part: part,
			rule: rule,
			Name: name,
			Mime: mime,
		}

		return
	}

	return
}
