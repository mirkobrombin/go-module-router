//go:build !ping
// +build !ping

package ping

type PingService interface{ Pong() string }

type pingService struct{}

func (pingService) Pong() string { return "pong" }
