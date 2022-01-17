package util

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/tylerb/gls"
)

type ridHook struct{}

func (h ridHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if id := gls.Get(ReqIDTag); id != nil && level != zerolog.NoLevel {
		e.Str(ReqIDTag, fmt.Sprintf("%v", id))
	}
}
