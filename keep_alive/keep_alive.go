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
	stopOnce  sync.Once

	log logrus.FieldLogger

	ticker   *time.Ticker
	stopChan chan struct{}

	service  KeepAliveService
	interval time.Duration

	maxFail int
}

type keepAliveOption func(a *AliveKeeper)

func NewKeepAlive(s KeepAliveService, opts ...keepAliveOption) *AliveKeeper {
	ak := &AliveKeeper{
		log:      logrus.New(),
		stopChan: make(chan struct{}),
		service:  s,
		interval: 5 * time.Second,
		maxFail:  3,
	}

	for _, opt := range opts {
		opt(ak)
	}
	ak.ticker = time.NewTicker(ak.interval)

	return ak
}

func KeepAlive(s KeepAliveService, opts ...keepAliveOption) *AliveKeeper {
	ak := NewKeepAlive(s, opts...)

	ak.Start()

	return ak
}
