package server

import (
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"sync"
)

type CmdHandler struct {
	itemsMux sync.RWMutex
	store    store.Store
}

func NewHandler() *CmdHandler {
	return &CmdHandler{
		store: store.GetStore(),
	}
}

func (h *CmdHandler) detach(conn Conn, cmd Command) {
	log := logger.WithField("cmd", cmd.Args[0])
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
	store.Submit(&common.TxRequest{
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
	val, err := h.store.Get(string(cmd.Args[1]))
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
	store.Submit(&common.TxRequest{
		TxId:   utils.GenerateId(),
		Flag:   common.FlagReq,
		Action: common.DEL,
		Key:    string(cmd.Args[1]),
	})

	h.itemsMux.Unlock()

	conn.WriteInt(1)
}
