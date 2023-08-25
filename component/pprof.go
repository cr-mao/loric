package component

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/cr-mao/loric/log"
)

var _ Component = &pprof{}

type pprof struct {
	Base
	addr string
}

func NewPProf(addr string) *pprof {
	return &pprof{
		addr: addr,
	}
}

func (p *pprof) Name() string {
	return "pprof"
}

func (p *pprof) Start() {
	go func() {
		log.Debug("pprof addr:", p.addr)
		err := http.ListenAndServe(p.addr, nil)
		if err != nil {
			log.Errorf("pprof server start failed: %v", err)
		}
	}()
}
