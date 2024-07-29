package keep_alive

import (
	"time"

	"github.com/sirupsen/logrus"
)

type KeepAliveService interface {
	GetName() string

	Start() error
	IsAlive() (bool, error)
	Kill() error
}

type AliveKeeper struct {
	log logrus.FieldLogger

	cancel <-chan struct{}

	service  KeepAliveService
	interval time.Duration
}

type keepAliveOption func(a *AliveKeeper)

func KeepAlive(s KeepAliveService, opts ...keepAliveOption) {
	ak := &AliveKeeper{
		log:      logrus.New(),
		cancel:   make(<-chan struct{}),
		service:  s,
		interval: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(ak)
	}

	ak.start()

	return
}

func (a *AliveKeeper) start() {
	t := time.NewTicker(a.interval)

	_ = a.startService()
	for {
		select {
		case <-t.C:
			if !a.startService() {
				return
			}
		case <-a.cancel:
			a.log.Infof(`Got termination call. Killing "%s"`, a.service.GetName())

			err := a.service.Kill()
			if err != nil {
				a.log.Errorf(`Error killing "%s"`, a.service.GetName())
			} else {
				a.log.Infof(`Successfully killed "%s"`, a.service.GetName())
			}

			return
		}
	}
}

func (a *AliveKeeper) startService() bool {
	isAlive, err := a.service.IsAlive()
	if err != nil {
		a.log.Errorf(`error checking if "%s" alive`, a.service.GetName())
		return false
	}
	if isAlive {
		return true
	}
	err = a.service.Kill()
	if err != nil {
		a.log.Errorf("error killing service %s", a.service.GetName())
		return false
	}
	err = a.service.Start()
	if err != nil {
		a.log.Errorf(`error keeping "%s" alive: %s `, a.service.GetName(), err)
		return false
	}
	a.log.Infof(`successfully started "%s"`, a.service.GetName())

	return true
}
