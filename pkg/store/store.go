package store

import (
	"encoding/json"
	"fmt"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"io/ioutil"
	"os"
	"path"
)

const (
	ERR = -(1 + iota)
	SET
	DEL
)

type Store interface {
	// Set new value, return old value if existed
	Set(key string, val []byte) (oldVal []byte, exist bool)
	// Get value of key
	Get(key string) (val []byte, err error)
	// Del value of key, and return value
	Del(key string) (val []byte, err error)
	// Serialize current data
	Serialize() (data []byte, err error)
	// Deserialize data to current state
	Load(data []byte) (err error)
}

type Evolvest struct {
	Nodes map[string][]byte `json:"nodes"`
}

var store Store

func init() {
	store = NewEvolvest()
}

func GetStore() Store {
	return store
}

func NewEvolvest() *Evolvest {
	return &Evolvest{Nodes: make(map[string][]byte, 17)}
}

func (e *Evolvest) Set(key string, val []byte) (oldVal []byte, exist bool) {
	oldVal, ok := e.Nodes[key]
	e.Nodes[key] = val

	defer func() {
		_ = GetWatcher().Notify(SET, key, oldVal, val)
	}()

	if ok {
		return oldVal, true
	}

	return nil, false
}

func (e *Evolvest) Get(key string) (val []byte, err error) {
	if val, ok := e.Nodes[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key %s not exists", key)
}

func (e *Evolvest) Del(key string) (val []byte, err error) {
	if val, ok := e.Nodes[key]; ok {

		delete(e.Nodes, key)
		_ = GetWatcher().Notify(DEL, key, val, nil)
		return val, nil
	}
	return nil, fmt.Errorf("key %s not exists", key)
}

func (e *Evolvest) Serialize() (data []byte, err error) {
	data, err = json.Marshal(e)
	return
}

func (e *Evolvest) Load(data []byte) (err error) {
	err = json.Unmarshal(data, e)
	return
}

func Persistent() {
	data, err := GetStore().Serialize()
	if err != nil {
		logger.Warn("save data error, %v", err)
		return
	}

	dataDir := config.Config().DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		logger.Warn("mkdir error, %v", err)
	}

	filename := path.Join(dataDir, common.SnapshotFile)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		logger.Warn("write data to file error, %v", err)
		return
	}
	logger.Info("write snapshot success!")
}

func Recover() {
	filename := path.Join(config.Config().DataDir, common.SnapshotFile)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warn("read data from file error, %v", err)
		return
	}
	err = GetStore().Load(data)
	if err != nil {
		logger.Warn("load data to store error, %v", err)
		return
	}

	logger.Info("recover data from snapshot success!")
}
