package logger

import (
	"github.com/go-logr/logr"
)

var (
	_ logr.LogSink = &Sink{}
)

// There should be no problem with concurent access and updates
type Sink struct {
	sink logr.LogSink
}

func NewSink(s logr.LogSink) *Sink {
	return &Sink{
		sink: s,
	}
}

func (s *Sink) SetSink(sink logr.LogSink) {
	s.sink = sink
}

func (s *Sink) Init(info logr.RuntimeInfo) {
	s.sink.Init(info)
}

func (d *Sink) Info(level int, msg string, keysAndValues ...interface{}) {
	d.sink.Info(level, msg, keysAndValues...)
}

func (d *Sink) Error(err error, msg string, keysAndValues ...interface{}) {
	d.sink.Error(err, msg, keysAndValues...)
}

func (d *Sink) Enabled(level int) bool {
	return d.sink.Enabled(level)
}

func (d *Sink) WithName(name string) logr.LogSink {
	return NewSink(d.sink.WithName(name))
}

func (d *Sink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return NewSink(d.sink.WithValues(keysAndValues...))
}
