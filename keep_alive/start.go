package keep_alive

import (
	"fmt"
	"time"

	"go.redsock.ru/rerrors"
)

// Start blocks until service become alive for the first time
func (a *AliveKeeper) Start() {
	a.startOnce.Do(func() {
		a.run()
		go func() {
			for range a.ticker.C {
				a.run()
				a.ticker.Reset(a.interval)
			}
			close(a.stopChan)
		}()

		return
	})
}

func (a *AliveKeeper) Stop() {
	a.stopOnce.Do(func() {
		a.ticker.Stop()
		<-a.stopChan
	})
}

func (a *AliveKeeper) run() {
	if a.service.IsAlive() {
		return
	}

	err := a.service.Kill()
	if err != nil {
		a.log.Error(rerrors.Wrap(err,
			fmt.Sprintf("error killing service %s", a.service.GetName())))
		return
	}
	err = a.service.Start()
	if err != nil {
		a.log.Error(rerrors.Wrap(err,
			fmt.Sprintf(`error keeping service %s alive`, a.service.GetName())))
		return
	}

	for retriesLeft := a.maxFail; retriesLeft > 0; retriesLeft-- {
		if a.service.IsAlive() {
			return
		}
		time.Sleep(a.interval)
	}

	a.log.Error(rerrors.New("failed healthcheks"))
}
