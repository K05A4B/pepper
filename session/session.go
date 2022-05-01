package session

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Session struct {
	Value          map[interface{}]interface{}
	LifeCycleStart int64

	change bool
}

// 设置 session 键值对
func (s *Session) Set(key interface{}, value interface{}) {
	s.change = true
	s.Value[key] = value
}

// 通过 key 获取数据
func (s *Session) Get(key interface{}) interface{} {
	data, ok := s.Value[key]

	if !ok {
		return nil
	}

	return data
}

// 删除一个键值对
func (s *Session) Remove(key interface{}) {
	s.change = true
	delete(s.Value, key)
}

// 清除对象内所有数据
func (s *Session) Empty() {
	s.change = true
	for k := range s.Value {
		delete(s.Value, k)
	}
}

// 判断有没有这对键值对
func (s *Session) Exist(key interface{}) bool {
	_, ok := s.Value[key]
	return ok
}

// 转成字符串数据
func (s *Session) ToBinary() (bin []byte, err error) {
	bytesBuf := new(bytes.Buffer)

	encoder := gob.NewEncoder(bytesBuf)
	err = encoder.Encode(s.Value)

	bin = bytesBuf.Bytes()

	return
}

// 创建 session 对象
// 参数为 nil 时反回 非 nil 的无数据 session 对象
// 参数为字符串指针时会尝试解析指针数据
func newSession() *Session {
	return &Session{
		Value:          make(map[interface{}]interface{}),
		LifeCycleStart: time.Now().Unix(),
	}
}
