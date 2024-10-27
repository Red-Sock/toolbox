package keep_alive

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type KeepAliveService interface {
	GetName() string

	Start() error
	IsAlive() bool
	Kill() error
}

type AliveKeeper struct {
	startOnce sync.Once

	log logrus.FieldLogger
	err error

	cancel <-chan struct{}

	service  KeepAliveService
	interval time.Duration

	maxFail int

	firstStartWG sync.WaitGroup
}

type keepAliveOption func(a *AliveKeeper)

func KeepAlive(s KeepAliveService, opts ...keepAliveOption) *AliveKeeper {
	ak := &AliveKeeper{
		log:      logrus.New(),
		cancel:   make(<-chan struct{}),
		service:  s,
		interval: 5 * time.Second,
		maxFail:  3,
	}

	for _, opt := range opts {
		opt(ak)
	}

	ak.firstStartWG.Add(1)
	ak.start()

	return ak
}

func (a *AliveKeeper) Wait() {
	a.firstStartWG.Wait()
}
