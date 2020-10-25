package store

import (
	"encoding/json"
	"fmt"
)

type Command interface {
	// Set new value, return old value if existed
	Set(key string, val string) (oldVal string, exist bool)
	// Get value of key
	Get(key string) (val string, err error)
	// Del value of key, and return value
	Del(key string) (val string, err error)
	// Save current data
	Save() (data []byte, err error)
	// Load data to current state
	Load(data []byte) (err error)
}

type Evolvest struct {
	Nodes map[string]string `json:"nodes"`
}

func NewEvolvest() *Evolvest {
	return &Evolvest{Nodes: make(map[string]string, 17)}
}

func (e *Evolvest) Set(key string, val string) (oldVal string, exist bool) {
	oldVal, ok := e.Nodes[key]
	e.Nodes[key] = val
	if ok {
		return oldVal, true
	}
	return "", false
}

func (e *Evolvest) Get(key string) (val string, err error) {
	if val, ok := e.Nodes[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key %s not exists", key)
}

func (e *Evolvest) Del(key string) (val string, err error) {
	if val, ok := e.Nodes[key]; ok {
		delete(e.Nodes, key)
		return val, nil
	}
	return "", fmt.Errorf("key %s not exists", key)
}

func (e *Evolvest) Save() (data []byte, err error) {
	data, err = json.Marshal(e)
	return
}

func (e *Evolvest) Load(data []byte) (err error) {
	err = json.Unmarshal(data, e)
	return
}
