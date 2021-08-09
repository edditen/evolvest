package server

import (
	"github.com/EdgarTeng/etlog"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"sync"
)

type CmdHandler struct {
	itemsMux sync.RWMutex
	syncer   *store.Syncer
}

func NewHandler(syncer *store.Syncer) *CmdHandler {
	return &CmdHandler{
		syncer: syncer,
	}
}

func (h *CmdHandler) detach(conn Conn, cmd Command) {
	log := etlog.Log.WithField("cmd", cmd.Args[0])
	detachedConn := conn.Detach()
	log.Info("connection has been detached")
	go func(c DetachedConn) {
		defer c.Close()

		c.WriteString("OK")
		c.Flush()
	}(detachedConn)
}

func (h *CmdHandler) ping(conn Conn, cmd Command) {
	conn.WriteString("PONG")
}

func (h *CmdHandler) quit(conn Conn, cmd Command) {
	conn.WriteString("OK")
	conn.Close()
}

func (h *CmdHandler) set(conn Conn, cmd Command) {
	if len(cmd.Args) != 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	h.itemsMux.Lock()
	h.syncer.Submit(&common.TxRequest{
		TxId:   utils.GenerateId(),
		Flag:   common.FlagReq,
		Action: common.SET,
		Key:    string(cmd.Args[1]),
		Val:    cmd.Args[2],
	})
	h.itemsMux.Unlock()

	conn.WriteString("OK")
}

func (h *CmdHandler) get(conn Conn, cmd Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	h.itemsMux.RLock()
	val, err := h.syncer.Store.Get(string(cmd.Args[1]))
	h.itemsMux.RUnlock()

	if err != nil {
		conn.WriteNull()
	} else {
		conn.WriteBulk(val.Val)
	}
}

func (h *CmdHandler) delete(conn Conn, cmd Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	h.itemsMux.Lock()
	h.syncer.Submit(&common.TxRequest{
		TxId:   utils.GenerateId(),
		Flag:   common.FlagReq,
		Action: common.DEL,
		Key:    string(cmd.Args[1]),
	})

	h.itemsMux.Unlock()

	conn.WriteInt(1)
}
