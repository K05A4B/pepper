package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/K05A4B/pepper/internal/utils"
)

type FileInfo struct {
	Size int64
	Name string
	Path string
	Mime string

	NotTrustworthyMime string
}

type Files struct {
	reader  *multipart.Reader
	rule    *Rule
	key     string
	partNum int
}

func (f *Files) Receive(dir string, size int) (info *FileInfo, err error) {
	var part *multipart.Part

	info = &FileInfo{}

	if f == nil {
		err = errors.New("invalid memory address or nil pointer dereference")
		return
	}

	for {
		part, err = f.reader.NextPart()
		if err != nil {
			return
		}

		info.Name = part.FileName()

		form := part.FormName()

		if form == f.key {
			break
		}
	}

	if (f.partNum+1) > f.rule.MaxNumber && f.rule.MaxNumber != 0 {
		err = errors.New("maximum number exceeded")
		return
	}

	info, err = f.save(part, dir, size)
	return
}

func (f *Files) save(part *multipart.Part, dir string, size int) (info *FileInfo, err error) {

	rule := f.rule

	filePath := dir + "/" + utils.GetRandString(40)
	if part == nil {
		err = errors.New("file is nil")
		return
	}

	info = &FileInfo{
		Name: part.FileName(),
		Path: filePath,
	}

	info = &FileInfo{
		Path: filePath,
		Name: part.FileName(),
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
	if err != nil {
		return
	}

	defer part.Close()
	defer file.Close()

	buffer := make([]byte, size)

	n := size
	for {
		n, err = part.Read(buffer[:n])
		if err != nil && err != io.EOF {
			return
		}

		if info.Mime == "" || info.NotTrustworthyMime == "" {
			info.Mime = http.DetectContentType(buffer[:n])
			info.NotTrustworthyMime = part.Header.Get("Content-Type")
		}

		if info.Size == 0 {
			// 检查 mime 是否允许
			if !rule.Mime.Exist(info.Mime) && rule.Mime.Len() != 0 {
				err = errors.New("this file type is not supported")
				return
			}
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			return
		}

		info.Size += int64(n)

		// 如果有规则信息则进行合格检测
		if rule != nil {
			if err == io.EOF || n == 0 && info.Size < rule.MinSize {
				file.Close()
				err = os.Remove(info.Path)
				err = errors.New("the file size is too small")
				part.Close()
				return
			}

			if info.Size > rule.MaxSize && rule.MaxSize > rule.MinSize {
				file.Close()
				err = os.Remove(info.Path)
				err = errors.New("the file is too large")
				part.Close()
				return
			}
		}

		if err == io.EOF || n == 0 {
			f.partNum++
			break
		}
	}

	return
}
