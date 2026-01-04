package logger

type nop struct{}

func (nop) Debug(msg string, kv ...any) {}
func (nop) Info(msg string, kv ...any)  {}
func (nop) Warn(msg string, kv ...any)  {}
func (nop) Error(msg string, kv ...any) {}
func (nop) Fatal(msg string, kv ...any) {}

var Nop Logger = nop{}
