package pkg

import "log"

type Logger interface {
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}

type Log struct{}

func (l Log) Infof(format string, a ...any) {
	log.Printf(format+"\n", a...)
}

func (l Log) Warnf(format string, a ...any) {
	log.Printf(format+"\n", a...)
}

func (l Log) Errorf(format string, a ...any) {
	log.Printf(format+"\n", a...)
}
