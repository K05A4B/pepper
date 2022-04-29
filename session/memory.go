package session

import (
	"sync"
	"time"
)

type Memory struct {
	Data map[string]*Session

	opt         *Options
	rwlock      sync.RWMutex
}

// 创建 session 对象
func (m *Memory) Create() (string, *Session) {
	id := getRandomSessionId()
	return id, m.Set(id)
}

// 删除 session 数据
func (m *Memory) Remove(id string) {
	delete(m.Data, id)
}

// 创建自定义 session id 的 session 对象
func (m *Memory) Set(id string) *Session {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	sess := newSession()
	m.Data[id] = sess

	time.AfterFunc(time.Duration(m.opt.LifeCycle)*time.Second, func() {
		id := id
		m.Remove(id)
	})

	return sess
}

// 获取 session 对象
func (m *Memory) Get(id string) *Session {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	sess, ok := m.Data[id]

	if !ok {
		return nil
	}

	if time.Now().Unix() > (sess.LifeCycleStart+m.opt.LifeCycle) {	
		m.Remove(id)	
	}

	return sess
}

// 检查 session id 对应的 session 对象是否存在
func (m *Memory) Exist(id string) bool {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	_, ok := m.Data[id]
	return ok
}

// 清除所有 session 对象
func (m *Memory) Clear() {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	for k := range m.Data {
		delete(m.Data, k)
	}
}

// 创建通过 内存来储存 Session 对象的管理器
func NewMemory(opt *Options) *Memory {
	if opt == nil {
		opt = &Options{
			LifeCycle: 60 * 60 * 24,
		}
	}

	return &Memory{
		Data: make(map[string]*Session),
		opt:  opt,
	}
}
