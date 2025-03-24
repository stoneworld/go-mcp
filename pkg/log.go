package pkg

import "log"

type Logger interface {
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}

type DefaultLogger struct{}

func (l DefaultLogger) Infof(format string, a ...any) {
	log.Printf(format+"\n", a...)
}

func (l DefaultLogger) Warnf(format string, a ...any) {
	log.Printf(format+"\n", a...)
}

func (l DefaultLogger) Errorf(format string, a ...any) {
	log.Printf(format+"\n", a...)
}
