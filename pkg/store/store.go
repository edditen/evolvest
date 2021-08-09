package store

import (
	"encoding/json"
	"fmt"
	"github.com/EdgarTeng/etlog"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/runnable"
	"io/ioutil"
	"os"
	"path"
)

type DataItem struct {
	Val []byte
	Ver int64
}

type Store interface {
	runnable.Runnable
	// Set new value, return old value if existed
	Set(key string, val DataItem) (oldVal DataItem, exist bool)
	// Get value of key
	Get(key string) (val DataItem, err error)
	// Del value of key, and return value
	Del(key string, ver int64) (val DataItem, err error)
	// Keys return all keys
	Keys() (keys []string, err error)
	// Serialize current data
	Serialize() (data []byte, err error)
	// Load data to current state
	Load(data []byte) (err error)
}

type Storage struct {
	cfg   *config.Config
	Nodes map[string]DataItem `json:"nodes"`
	w     *Watcher
}

func NewStorage(conf *config.Config) *Storage {
	return &Storage{
		cfg:   conf,
		w:     NewWatcher(),
		Nodes: make(map[string]DataItem, 17),
	}
}

func (s *Storage) Init() error {
	return nil
}

func (s *Storage) Run() error {
	return nil
}

func (s *Storage) Shutdown() {
}

func (s *Storage) Set(key string, val DataItem) (oldVal DataItem, exist bool) {
	oldVal, ok := s.Nodes[key]
	if ok && val.Ver < oldVal.Ver {
		// exist key, compare with the original one
		return oldVal, true
	}
	s.Nodes[key] = val

	defer func() {
		_ = s.w.Notify(common.SET, key, oldVal, val)
	}()

	if ok {
		return oldVal, true
	}

	return DataItem{}, false
}

func (s *Storage) Get(key string) (val DataItem, err error) {
	if val, ok := s.Nodes[key]; ok {
		return val, nil
	}
	return DataItem{}, fmt.Errorf("key %s not exists", key)
}

func (s *Storage) Del(key string, ver int64) (val DataItem, err error) {
	if val, ok := s.Nodes[key]; ok {
		if ver < val.Ver {
			return DataItem{}, fmt.Errorf("ver %d is less than Store", ver)
		}
		delete(s.Nodes, key)
		_ = s.w.Notify(common.DEL, key, val, DataItem{})
		return val, nil
	}
	return DataItem{}, fmt.Errorf("key %s not exists", key)
}

func (s *Storage) Keys() (keys []string, err error) {
	keys = make([]string, 0, len(s.Nodes))
	for k := range s.Nodes {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *Storage) Serialize() (data []byte, err error) {
	data, err = json.Marshal(s)
	return
}

func (s *Storage) Load(data []byte) (err error) {
	err = json.Unmarshal(data, s)
	return
}

func (s *Storage) Persistent() {
	data, err := s.Serialize()
	if err != nil {
		etlog.Log.WithError(err).Warn("save data error")
		return
	}

	dataDir := s.cfg.DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		etlog.Log.WithError(err).Warn("mkdir error")
	}

	filename := path.Join(dataDir, common.FileSnapshot)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		etlog.Log.WithError(err).Warn("write data to file error")
		return
	}
	etlog.Log.Info("write snapshot success!")
}

func (s *Storage) Recover() {
	filename := path.Join(s.cfg.DataDir, common.FileSnapshot)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		etlog.Log.WithError(err).Warn("read data from file error")
		return
	}
	err = s.Load(data)
	if err != nil {
		etlog.Log.WithError(err).Warn("load data to Store error")
		return
	}

	etlog.Log.Info("recover data from snapshot success!")
}
