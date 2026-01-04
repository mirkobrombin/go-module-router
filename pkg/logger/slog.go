package logger

import "log/slog"

type SlogLogger struct {
	l *slog.Logger
}

func NewSlog(l *slog.Logger) Logger {
	return &SlogLogger{l: l}
}

func (s *SlogLogger) Debug(msg string, kv ...any) { s.l.Debug(msg, kv...) }
func (s *SlogLogger) Info(msg string, kv ...any)  { s.l.Info(msg, kv...) }
func (s *SlogLogger) Warn(msg string, kv ...any)  { s.l.Warn(msg, kv...) }
func (s *SlogLogger) Error(msg string, kv ...any) { s.l.Error(msg, kv...) }
func (s *SlogLogger) Fatal(msg string, kv ...any) { s.l.Error(msg, kv...) }
