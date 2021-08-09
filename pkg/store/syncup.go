package store

import (
	"errors"
	"github.com/edditen/evolvest/pkg/common"
	"github.com/edditen/evolvest/pkg/common/config"
)

type Syncer struct {
	cfg      *config.Config
	Store    Store
	appender Appender
	sender   Sender
	reqC     chan *common.TxRequest
	shutdown chan interface{}
}

func NewSyncer(conf *config.Config) *Syncer {
	return &Syncer{
		cfg:      conf,
		Store:    NewStorage(conf),
		appender: NewTxAppender(conf),
		sender:   NewTxSender(conf),
		reqC:     make(chan *common.TxRequest, 1000),
		shutdown: make(chan interface{}),
	}
}

func (s *Syncer) Init() error {
	if err := s.Store.Init(); err != nil {
		return err
	}
	if err := s.appender.Init(); err != nil {
		return err
	}
	if err := s.sender.Init(); err != nil {
		return err
	}
	return nil
}

func (s *Syncer) Run() error {
	defer close(s.reqC)
	for {
		select {
		case req := <-s.reqC:
			s.setToStore(req)
			s.appender.Append(req)
			if req.Flag == common.FlagReq {
				go s.sender.Send(req)
			}
		case <-s.shutdown:
			break
		}

	}

	return nil

}

func (s *Syncer) Shutdown() {
	s.sender.Shutdown()
	s.appender.Shutdown()
	s.Store.Shutdown()
	close(s.shutdown)
}

func (s *Syncer) Submit(req *common.TxRequest) error {
	select {
	case s.reqC <- req:
		return nil
	default:
		return errors.New("tx chan is full or off")
	}
}

func (s *Syncer) setToStore(req *common.TxRequest) {
	switch req.Action {
	case common.SET:
		s.Store.Set(req.Key, DataItem{
			Val: req.Val,
			Ver: req.TxId,
		})
	case common.DEL:
		s.Store.Del(req.Key, req.TxId)
	}
}
