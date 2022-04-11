package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"os"

	"github.com/kz91/pepper/internal/utils"
)

type File struct {
	Size int64
	Name string
	Mime string
	Path string

	part *multipart.Part
	rule *Rule
	file *os.File
}

func (f *File) Receive(dir string, bufferSize int) (fop *File, err error) {
	rule := f.rule
	part := f.part
	if part == nil {
		err = errors.New("file is nil")
		return
	}

	filePath := dir + "/" + utils.GetRandString(40)
	f.Path = filePath

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
	if err != nil {
		return
	}

	f.file = file

	defer part.Close()
	defer file.Close()

	buffer := make([]byte, bufferSize)
	var fileSize int64 = 0

	for {
		n := bufferSize

		n, err = part.Read(buffer[:n])
		if err != nil && err != io.EOF {
			return
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			return
		}

		fileSize += int64(n)

		// 如果有规则信息则进行合格检测
		if rule != nil {
			if err == io.EOF || n == 0 && fileSize < rule.MinSize {
				err = f.Remove()
				err = errors.New("the file size is too small")
				part.Close()
				return
			}

			if fileSize > rule.MaxSize && rule.MaxSize > rule.MinSize {
				err = f.Remove()
				err = errors.New("the file is too large")
				part.Close()
				return
			}
		}

		if err == io.EOF || n == 0 {
			break
		}
	}

	fop = f
	return
}

func (f *File) Remove() error {
	p := f.Path
	f.Path = ""

	if f.file != nil {
		f.file.Close()
	}

	if err := os.Remove(p); err != nil {
		return err
	}
	return nil
}
