package session

import (
	"sync"
	"time"
)

type Memory struct {
	Data map[string]*Session

	opt    *Options
	rwlock sync.RWMutex
}

// 创建 session 对象
func (m *Memory) Create() (string, *Session) {
	id := getRandomSessionId()
	return id, m.set(id)
}

// 删除 session 数据
func (m *Memory) Remove(id string) {
	delete(m.Data, id)
}

// 获取 session 对象
func (m *Memory) Get(id string) *Session {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	sess, ok := m.Data[id]

	if !ok {
		return nil
	}

	if time.Now().Unix() > (sess.LifeCycleStart + m.opt.LifeCycle) {
		m.Remove(id)
		return nil
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

// 清理生命周期结束的 session 对象
func (m *Memory) Clean() {
	for id, sess := range m.Data {
		if time.Now().Unix() > (sess.LifeCycleStart + m.opt.LifeCycle) {
			m.Remove(id)
		}
	}
}

// 清除所有 session 对象
func (m *Memory) Empty() {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	for k := range m.Data {
		delete(m.Data, k)
	}
}

// 创建自定义 session id 的 session 对象
func (m *Memory) set(id string) *Session {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	sess := newSession()
	m.Data[id] = sess

	return sess
}

// 初始化options
func (m *Memory) init() {
	if m.opt.CleanInterval == 0 {
		m.opt.CleanInterval = 60 * 60 * 12
	}

	if m.opt.LifeCycle == 0 {
		m.opt.LifeCycle = 60 * 60 * 24
	}

	if m.opt.HandlerGCError == nil {
		m.opt.HandlerGCError = func(err error) {}
	}
}

func (m *Memory) gc() {
	time.AfterFunc(time.Duration(m.opt.CleanInterval)*time.Second, func() {
		m.Clean()
		m.gc()
	})
}

// 创建通过 内存来储存 Session 对象的管理器
func NewMemory(opt *Options) *Memory {
	if opt == nil {
		opt = &Options{
			LifeCycle: 60 * 60 * 24,
		}
	}

	m := &Memory{
		Data: make(map[string]*Session),
		opt:  opt,
	}

	m.init()
	m.gc()

	return m
}
