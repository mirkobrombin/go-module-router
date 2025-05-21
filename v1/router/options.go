package router

import (
	"time"

	"github.com/mirkobrombin/go-module-router/v1/logger"
)

type Options struct {
	SessionDuration time.Duration
	Logger          logger.Logger
	OnError         func(error)
	SkipAutoWire    bool
}
