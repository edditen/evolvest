package store

import "fmt"

type Command interface {
	// Set new value, return old value if existed
	Set(key string, val string) (oldVal string, exist bool)
	// Get value of key
	Get(key string) (val string, err error)
	// Del value of key, and return value
	Del(key string) (val string, err error)
}

type Evolvest struct {
	storage map[string]string
}

func NewEvolvest() *Evolvest {
	return &Evolvest{storage: make(map[string]string, 17)}
}

func (e *Evolvest) Set(key string, val string) (oldVal string, exist bool) {
	oldVal, ok := e.storage[key]
	e.storage[key] = val
	if ok {
		return oldVal, true
	} else {
		return "", false
	}
}

func (e *Evolvest) Get(key string) (val string, err error) {
	if val, ok := e.storage[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key %s not exists", key)
	}
}

func (e *Evolvest) Del(key string) (val string, err error) {
	if val, ok := e.storage[key]; ok {
		delete(e.storage, key)
		return val, nil
	} else {
		return "", fmt.Errorf("key %s not exists", key)
	}
}
