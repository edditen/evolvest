package store

import (
	"fmt"
	"github.com/edditen/etlog"
	"github.com/edditen/evolvest/pkg/common"
	"github.com/edditen/evolvest/pkg/common/config"
	"github.com/edditen/evolvest/pkg/common/utils"
	"github.com/edditen/evolvest/pkg/runnable"
	"github.com/pkg/errors"
	"log"
	"os"
	"path"
)

type Appender interface {
	runnable.Runnable
	Append(req *common.TxRequest) error
}

type TxAppender struct {
	cfg      *config.Config
	writer   *os.File
	shutdown chan interface{}
}

func NewTxAppender(cfg *config.Config) *TxAppender {
	return &TxAppender{
		cfg:      cfg,
		shutdown: make(chan interface{}),
	}
}

func (ta *TxAppender) Init() error {
	log.Println("[Init] init txAppender")
	dataDir := ta.cfg.DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "init syncUp error")
	}
	filename := path.Join(dataDir, common.FileTx)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open tx file failed")
	}
	ta.writer = f
	return nil
}

func (ta *TxAppender) Run(errC chan<- error) {
	log.Println("[Run] run txAppender")
}

func (ta *TxAppender) Shutdown() {
}

func (ta *TxAppender) Append(req *common.TxRequest) error {
	text := fmt.Sprintf("%d %s %s %s %s\n",
		req.TxId, req.Flag, req.Action, req.Key, utils.Base64Encode(req.Val))
	if _, err := ta.writer.WriteString(text); err != nil {
		etlog.Log.WithError(err).
			WithField("text", text).
			Warn("append text to tx file failed")
		return errors.Wrap(err, "append tx to file error")
	}
	return nil
}
