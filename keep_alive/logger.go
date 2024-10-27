package keep_alive

import (
	"github.com/sirupsen/logrus"
)

func WithoutLogger() keepAliveOption {
	return func(a *AliveKeeper) {
		l := logrus.New()
		l.SetOutput(&nullWriter{})
		a.log = l
	}
}

type nullWriter struct{}

func (nullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
