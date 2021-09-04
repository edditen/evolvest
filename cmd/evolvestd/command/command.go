package command

import (
	"github.com/edditen/etlog"
	"github.com/edditen/evolvest/embed/rpc"
	"github.com/edditen/evolvest/embed/server"
	"github.com/edditen/evolvest/pkg/common/config"
	"github.com/edditen/evolvest/pkg/store"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
)

type Evolvestd struct {
	config         *config.Config
	syncer         *store.Syncer
	syncServer     *rpc.SyncServer
	evolvestServer *server.EvolvestServer
}

func NewEvolvestd() *Evolvestd {
	return &Evolvestd{}
}

func (e *Evolvestd) Init() (err error) {
	logger, err := etlog.NewEtLogger(etlog.SetConfigPath(viper.GetString("log-config")))
	if err != nil {
		return errors.Wrap(err, "init logger error")
	}
	etlog.SetDefaultLog(logger)

	e.config = config.NewConfig(viper.GetString("config"))
	if err = e.config.Init(); err != nil {
		return errors.Wrap(err, "init config error")
	}

	e.syncer = store.NewSyncer(e.config)
	if err = e.syncer.Init(); err != nil {
		return errors.Wrap(err, "init syncer error")
	}

	e.syncServer = rpc.NewSyncServer(e.config, e.syncer)
	if err = e.syncServer.Init(); err != nil {
		return errors.Wrap(err, "init syncServer error")
	}

	e.evolvestServer = server.NewEvolvestServer(e.config, e.syncer)
	if err = e.evolvestServer.Init(); err != nil {
		return errors.Wrap(err, "init evolvestServer error")
	}

	return nil
}

func (e *Evolvestd) Run(errC chan<- error) {
	go e.config.Run(errC)
	go e.syncer.Run(errC)
	go e.syncServer.Run(errC)
	go e.evolvestServer.Run(errC)
}

func (e *Evolvestd) Shutdown() {
	e.evolvestServer.Shutdown()
	e.syncServer.Shutdown()
	e.syncer.Shutdown()
	e.config.Shutdown()
}

func (e *Evolvestd) WaitSignal(errC <-chan error, hook func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	select {
	case <-c:
		log.Println("Server received interrupt signal")
		hook()
		os.Exit(0)
	case err := <-errC:
		log.Printf("Server run error: %+v", err)
		hook()
		os.Exit(1)
	}
}
