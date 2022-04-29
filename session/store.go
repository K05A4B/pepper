package session

import (
	"encoding/gob"
	"io"
	"os"
	"sync"
	"time"
)

type Store struct {
	Memory
	file string
	fp   *os.File
	lock sync.Mutex
}

// 重新加载数据
func (s *Store) Reload() error {
	var err error

	s.lock.Lock()
	defer s.lock.Unlock()

	decoder := gob.NewDecoder(s.fp)

	err = decoder.Decode(&s.Data)
	if err != nil {
		return err
	}

	
	for id, sess := range s.Data {
		if time.Now().Unix() > (sess.LifeCycleStart + s.opt.LifeCycle) {
			s.Remove(id)
			s.Save()
		}
	}

	return nil
}

// 保存数据到文件
func (s *Store) Save() error {
	encoder := gob.NewEncoder(s.fp)

	err := s.fp.Truncate(0)
	if err != nil {
		return err
	}

	s.fp.Seek(0, 0)

	err = encoder.Encode(s.Data)
	if err != io.EOF {
		return err
	}

	return nil
}

// 关闭对文件的操作
func (s *Store) Close() error {
	return s.fp.Close()
}

// 加载数据文件
func (s *Store) load() error {
	var err error

	if s.fp == nil {
		s.fp, err = os.OpenFile(s.file, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
	}

	return s.Reload()
}

// 创建通过 文件加内存来储存 Session 对象的管理器
func NewStore(file string, opt *Options) (*Store, error) {
	f := &Store{
		file:   file,
		Memory: *NewMemory(opt),
	}

	err := f.load()
	if err == io.EOF {
		return f, nil
	}

	return f, err
}
