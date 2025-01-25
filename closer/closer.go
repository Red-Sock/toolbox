package closer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sync"
)

type Closer struct {
	m         sync.Mutex
	isClosed  bool
	scheduled []*closable
}

type Closable func() error

type closable struct {
	Close    Closable
	isClosed bool
	err      error
}

func (c *Closer) Add(f Closable) {
	c.m.Lock()
	c.scheduled = append(c.scheduled, &closable{Close: f})
	c.m.Unlock()
}

func (c *Closer) Close() (err error) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, f := range c.scheduled {
		f.err = f.Close()
		if f.err != nil {
			err = errors.Join(err, f.err)
		} else {
			f.isClosed = true
		}
	}

	return err
}

type State struct {
	Name     string `json:"name"`
	IsClosed bool   `json:"is_closed"`
	Error    string `json:"error"`
}

func (c *Closer) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	c.m.Lock()
	defer c.m.Unlock()
	states := make([]State, 0, len(c.scheduled))

	for _, s := range c.scheduled {
		pc := runtime.FuncForPC(reflect.ValueOf(s.Close).Pointer())
		st := State{
			Name:     fmt.Sprintf("%s", pc.Name()),
			IsClosed: s.isClosed,
		}

		if s.err != nil {
			st.Error = s.err.Error()
		}

		states = append(states, st)
	}

	_ = json.NewEncoder(w).Encode(states)
}
