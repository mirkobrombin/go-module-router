package logger

import "go.uber.org/zap"

type Zap struct{ L *zap.Logger }

func (z *Zap) Debug(msg string, kv ...any) { z.L.Sugar().Debugw(msg, kv...) }
func (z *Zap) Info(msg string, kv ...any)  { z.L.Sugar().Infow(msg, kv...) }
func (z *Zap) Warn(msg string, kv ...any)  { z.L.Sugar().Warnw(msg, kv...) }
func (z *Zap) Error(msg string, kv ...any) { z.L.Sugar().Errorw(msg, kv...) }
func (z *Zap) Fatal(msg string, kv ...any) { z.L.Sugar().Fatalw(msg, kv...) }
