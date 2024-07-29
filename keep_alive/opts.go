package keep_alive

import (
	"time"

	"github.com/sirupsen/logrus"
)

func WithLogger(logger logrus.FieldLogger) keepAliveOption {
	return func(a *AliveKeeper) {
		a.log = logger
	}
}

func WithCheckInterval(interval time.Duration) keepAliveOption {
	return func(a *AliveKeeper) {
		a.interval = interval
	}
}
func WithCancel(c <-chan struct{}) keepAliveOption {
	return func(a *AliveKeeper) {
		a.cancel = c
	}
}
