package keep_alive

import (
	"fmt"
	"sync"
	"time"

	"go.redsock.ru/rerrors"
)

// Start blocks until service become alive for the first time
func (a *AliveKeeper) Start() {
	a.startOnce.Do(a.start)
}

func (a *AliveKeeper) start() {
	t := time.NewTicker(time.Nanosecond)

	failedTimes := a.maxFail
	onceDone := sync.Once{}

	go func() {
		for {
			select {
			case <-t.C:
				t.Reset(a.interval)

				startErr := a.startService()
				if startErr == nil {
					onceDone.Do(a.firstStartWG.Done)
					failedTimes = a.maxFail
					continue
				}

				failedTimes--
				if failedTimes != 0 {
					a.log.Errorf("error keeping alive service %s. Error: %s",
						a.service.GetName(), startErr)
					continue
				}

				a.log.Errorf("keep alive failed %d times. Last error: %s",
					failedTimes, startErr)
				return
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
	}()

	return
}

func (a *AliveKeeper) startService() error {
	if a.service.IsAlive() {
		return nil
	}

	err := a.service.Kill()
	if err != nil {
		return rerrors.Wrap(err,
			fmt.Sprintf("error killing service %s", a.service.GetName()))
	}
	err = a.service.Start()
	if err != nil {
		return rerrors.Wrap(err,
			fmt.Sprintf(`error keeping service %s alive`, a.service.GetName()))
	}

	for retriesLeft := a.maxFail; retriesLeft > 0; retriesLeft-- {
		if a.service.IsAlive() {
			return nil
		}
		time.Sleep(a.interval)
	}

	return rerrors.New("failed healthcheks")
}
