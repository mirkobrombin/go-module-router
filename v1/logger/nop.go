package logger

type nop struct{}

func (nop) Debug(string, ...any) {}
func (nop) Info(string, ...any)  {}
func (nop) Warn(string, ...any)  {}
func (nop) Error(string, ...any) {}
func (nop) Fatal(string, ...any) {}

var Nop Logger = nop{}
