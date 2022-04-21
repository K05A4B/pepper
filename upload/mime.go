package upload

import "strings"

type Mime struct {
	Mime []string
}

func (r *Mime) Find(m string) int {
	for index, mime := range r.Mime {
		if mime == m {
			return index
		}
	}

	return -1
}

func (m *Mime) Append(s string) {
	m.Mime = append(m.Mime, s)
}

func (m *Mime) Exist(s string) bool {
	s = strings.Split(s, ";")[0]
	if m.Find(s) < 0 {
		return false
	} else {
		return true
	}
}

func (m *Mime) Len() int {
	return len(m.Mime)
}
