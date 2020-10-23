package store

import "fmt"

type Command interface {
	// Set new value, return old value if existed
	Set(key string, val interface{}) (oldVal interface{}, err error)
	// Get value of key
	Get(key string) (val interface{}, err error)
	// Del value of key, and return value
	Del(key string) (val interface{}, err error)
}

type Evolvest struct {
	storage map[string]interface{}
}

func NewEvolvest() *Evolvest {
	return &Evolvest{storage: make(map[string]interface{}, 17)}
}

func (e *Evolvest) Set(key string, val interface{}) (oldVal interface{}, err error) {
	oldVal, _ = e.storage[key]
	e.storage[key] = val
	return oldVal, nil
}

func (e *Evolvest) Get(key string) (val interface{}, err error) {
	if val, ok := e.storage[key]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("key %s not exists", key)
	}
}

func (e *Evolvest) Del(key string) (val interface{}, err error) {
	if val, ok := e.storage[key]; ok {
		delete(e.storage, key)
		return val, nil
	} else {
		return nil, nil
	}
}
