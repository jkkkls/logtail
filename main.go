package main

import (
	"log"
	"syscall"

	"github.com/judwhite/go-svc"
)

type LogTail struct{}

func main() {
	logtail := &LogTail{}
	if err := svc.Run(logtail, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatalf("svc.Run() returned: %s", err)
	}
}

func (l *LogTail) Init(env svc.Environment) error {
	return nil
}

func (l *LogTail) Start() error {
	return nil
}

func (l *LogTail) Stop() error {
	return nil
}
