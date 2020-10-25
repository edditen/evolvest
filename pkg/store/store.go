package store

import (
	"encoding/json"
	"fmt"
)

const (
	ERR = -(1 + iota)
	SET
	DEL
)

type Store interface {
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

var evolvest *Evolvest

func init() {
	evolvest = NewEvolvest()
}

func GetEvolvest() *Evolvest {
	return evolvest
}

func NewEvolvest() *Evolvest {
	return &Evolvest{Nodes: make(map[string]string, 17)}
}

func (e *Evolvest) Set(key string, val string) (oldVal string, exist bool) {
	oldVal, ok := e.Nodes[key]
	e.Nodes[key] = val

	defer func() {
		GetWatcher().Notify(SET, key, oldVal, val)
	}()

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
		GetWatcher().Notify(DEL, key, val, "")
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
