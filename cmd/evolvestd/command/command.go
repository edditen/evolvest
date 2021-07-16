package command

import (
	"github.com/EdgarTeng/etlog"
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/embed/server"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/store"
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

func (e *Evolvestd) Run() error {
	errC := make(chan error)

	go e.runConfig(errC)
	go e.runSyncer(errC)
	go e.runSyncServer(errC)
	go e.runEvolvestServer(errC)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Server received interrupt signal, prepare to clean.")
		e.evolvestServer.Shutdown()
		e.syncServer.Shutdown()
		e.syncer.Shutdown()

		log.Println("Server received interrupt signal, exit now.")
		os.Exit(0)
	}()

	return <-errC
}

func (e *Evolvestd) Shutdown() {
	return
}

func (e *Evolvestd) runConfig(errC chan error) {
	if err := e.config.Run(); err != nil {
		errC <- errors.Wrap(err, "run config error")
	}
}

func (e *Evolvestd) runSyncer(errC chan error) {
	if err := e.syncer.Run(); err != nil {
		errC <- errors.Wrap(err, "run syncer error")
	}
}

func (e *Evolvestd) runSyncServer(errC chan error) {
	if err := e.syncServer.Run(); err != nil {
		errC <- errors.Wrap(err, "run syncServer error")
	}
}

func (e *Evolvestd) runEvolvestServer(errC chan error) {
	if err := e.evolvestServer.Run(); err != nil {
		errC <- errors.Wrap(err, "run evolvestServer error")
	}
}
