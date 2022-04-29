package session

import "time"

type Session struct {
	Value          map[interface{}]interface{}
	LifeCycleStart int64
}

// 设置 session 键值对
func (s *Session) Set(key interface{}, value interface{}) {
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
	delete(s.Value, key)
}

// 清除对象内所有数据
func (s *Session) Clear() {
	for k := range s.Value {
		delete(s.Value, k)
	}
}

func (s *Session) Exist(key interface{}) bool {
	_, ok := s.Value[key]
	return ok
}

func newSession() *Session {
	return &Session{
		Value: make(map[interface{}]interface{}),
		LifeCycleStart: time.Now().Unix(),
	}
}
