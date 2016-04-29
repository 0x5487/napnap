package napnap

import "log"

const (
	debugLevel   = 0
	infoLevel    = 1
	warningLevel = 2
	errorLevel   = 3
	fatalLevel   = 4
	off          = 5
)

type logger struct {
	mode int
}

func newLog() *logger {
	return &logger{
		mode: infoLevel,
	}
}

func (l *logger) debug(v ...interface{}) {
	if l.mode <= debugLevel {
		log.Println(v)

	}
}

func (l *logger) debugf(format string, v ...interface{}) {
	if l.mode <= debugLevel {
		log.Printf(format, v)
	}
}

func (l *logger) fatal(v ...interface{}) {
	if l.mode <= fatalLevel {
		log.Fatal(v)
	}
}

func (l *logger) fatalf(format string, v ...interface{}) {
	if l.mode <= fatalLevel {
		log.Fatalf(format, v)
	}
}
