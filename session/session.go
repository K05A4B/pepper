package session

import (
	"encoding/gob"
	"errors"
	"io"
	"os"

	"github.com/kz91/pepper/internal/utils"
)

type Session map[string]interface{}

type Manager struct {
	Data map[string]Session
	file string
}

// 加载SESSION文件
func (m *Manager) Load(file string) error {
	m.file = file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDONLY, 0777)

	if err != nil {
		return err
	}

	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&m.Data)

	if err == io.EOF {
		return f.Close()
	}

	if err != nil {
		return err
	}

	return f.Close()
}

// 保存 session
func (m *Manager) Dump() error {
	if m.file == "" {
		return errors.New("file not found")
	}
	f, err := os.OpenFile(m.file, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(m.Data); err != nil {
		return err
	}

	return f.Close()
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
func (m *Manager) Get(id string) *Session {
	if m.Data == nil {
		m.Data = make(map[string]Session)
	}

	Data, ok := m.Data[id]
	if !ok {
		return nil
	}

	return &Data
}

func NewManager() *Manager {
	return &Manager{}
}