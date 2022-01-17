package util_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tylerb/gls"
	"math/rand"
	"signature-server/util"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("no id", func(t *testing.T) {
		util.Info("Hello World! no id")
		gls.Cleanup()
	})

	t.Run("info", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Info("Hello World!")
		gls.Cleanup()
	})

	t.Run("infof", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Infof("%s %s", "Hello", "World!")
		gls.Cleanup()
	})

	t.Run("debug", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Debug("Hello World!")
		gls.Cleanup()
	})

	t.Run("debug", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Debugf("%s %s", "Hello", "World!")
		gls.Cleanup()
	})

	t.Run("warn", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Warn("Hello World!")
		gls.Cleanup()
	})
	t.Run("warn", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Warnf("%s %s", "Hello", "World!")
		gls.Cleanup()
	})

	t.Run("error", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Error("Hello World!")
		gls.Cleanup()
	})

	t.Run("errorf", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		util.Errorf("%s %s", "Hello", "World!")
		gls.Cleanup()
	})

	t.Run("panic", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		assert.Panics(t, func() { util.Panic("Hello World!") })
		gls.Cleanup()
	})

	t.Run("panicf", func(t *testing.T) {
		gls.Set(util.ReqIDTag, rand.Int63())
		assert.Panics(t, func() { util.Panicf("%s %s", "Hello", "World!") })
		gls.Cleanup()
	})
}
