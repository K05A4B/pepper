package session

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"

	"github.com/kz91/pepper/internal/utils"
)

type Session map[string]interface{}

type Manager struct {
	Data map[string]Session
	file string
	fp   *os.File
	bfio *bufio.Writer
}

// 加载SESSION文件
func (m *Manager) Load(file string) error {
	m.file = file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil && err != io.EOF {
		return err
	}

	m.fp = f
	m.bfio = bufio.NewWriter(f)

	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&m.Data)

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

// 创建一个用户的session
func (m *Manager) Create() string {
	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	id := "SESSION_ID_" + utils.GetRandString(40)
	m.Data[id] = make(Session)
	return id
}

// 通过session id 获取session
func (m *Manager) Get(id string, key string) interface{} {
	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	data, ok := m.Data[id]
	if !ok {
		return nil
	}

	value, ok := data[key]
	if !ok {
		return nil
	}

	return value
}

func (m *Manager) Set(id string, key string, value interface{}) error {
	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	data := m.Data[id]
	if data == nil {
		data = make(Session)
	}
	data[key] = value

	return m.Flush()
}

func (m *Manager) Flush() error {
	encoder := gob.NewEncoder(m.fp)
	if err := encoder.Encode(m.Data); err != nil {
		return err
	}
	return m.bfio.Flush()
}

func (m *Manager) Close() error {
	return m.fp.Close()
}

func NewManager() *Manager {
	return &Manager{}
}
