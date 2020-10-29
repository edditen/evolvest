package store

import (
	"encoding/json"
	"fmt"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/log"
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
	Set(key string, val string) (oldVal string, exist bool)
	// Get value of key
	Get(key string) (val string, err error)
	// Del value of key, and return value
	Del(key string) (val string, err error)
	// Serialize current data
	Serialize() (data []byte, err error)
	// Deserialize data to current state
	Load(data []byte) (err error)
}

type Evolvest struct {
	Nodes map[string]string `json:"nodes"`
}

var store Store

func init() {
	store = NewEvolvest()
}

func GetStore() Store {
	return store
}

func NewEvolvest() *Evolvest {
	return &Evolvest{Nodes: make(map[string]string, 17)}
}

func (e *Evolvest) Set(key string, val string) (oldVal string, exist bool) {
	oldVal, ok := e.Nodes[key]
	e.Nodes[key] = val

	defer func() {
		_ = GetWatcher().Notify(SET, key, oldVal, val)
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
		_ = GetWatcher().Notify(DEL, key, val, "")
		return val, nil
	}
	return "", fmt.Errorf("key %s not exists", key)
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
		log.Warn("save data error, %v", err)
		return
	}

	dataDir := config.Config().DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Warn("mkdir error, %v", err)
	}

	filename := path.Join(dataDir, common.SnapshotFile)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Warn("write data to file error, %v", err)
		return
	}
	log.Info("write snapshot success!")
}

func Recover() {
	filename := path.Join(config.Config().DataDir, common.SnapshotFile)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Warn("read data from file error, %v", err)
		return
	}
	err = GetStore().Load(data)
	if err != nil {
		log.Warn("load data to store error, %v", err)
		return
	}

	log.Info("recover data from snapshot success!")
}
