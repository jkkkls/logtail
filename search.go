package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hpcloud/tail"
)

type LogConfig struct {
	ConfigName    string
	Config        *Config
	Offset        int64
	WatchFileName string
	Close         chan struct{}
	run           int32
	w             *sync.WaitGroup
	sinkWriter    SinkWriter
	parser        Parser
}

var logFiles = map[string]*LogConfig{}

func (l *LogConfig) reset(m map[string]int64) <-chan time.Time {
	now := time.Now()
	if strings.Contains(l.Config.File, "{date}") {
		fileName := strings.ReplaceAll(l.Config.File, "{date}", now.Format("2006-01-02"))
		l.WatchFileName = fileName
		l.Offset = m[l.WatchFileName]
		// 第二天0点0分0秒的定时器
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		timer := time.NewTimer(nextDay.Sub(now))
		return timer.C
	}

	l.WatchFileName = l.Config.File
	l.Offset = m[l.WatchFileName]
	return make(<-chan time.Time, 1)
}

func (l *LogConfig) Run(m map[string]int64) {
	c := l.reset(m)

	if !atomic.CompareAndSwapInt32(&l.run, 0, 1) {
		return
	}
	defer func() {
		atomic.StoreInt32(&l.run, 0)
	}()

	textLogger.Info("tail file", "file", l.WatchFileName)
	if l.WatchFileName == "" {
		<-l.Close
		l.w.Done()
		return
	}

	t, err := tail.TailFile(l.WatchFileName, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: l.Offset}})
	if err != nil {
		textLogger.Warn("tail file error: ", "err", err, "file", l.WatchFileName)
		return
	}

	for {
		select {
		case line := <-t.Lines:
			// add \n length
			l.Offset += int64(len(line.Text)) + 1
			parsed, err := l.parser.Parse(line.Text)
			if err != nil {
				textLogger.Warn("parse text error", "err", err)
				continue
			}

			err = l.sinkWriter.Write(parsed)
			if err != nil {
				textLogger.Warn("write to sink error", "err", err)
			}
		case <-l.Close:
			if l.w != nil {
				l.w.Done()
			}
			return
		case <-c:
			// 第二天
			c = l.reset(map[string]int64{})

			t, err = tail.TailFile(l.WatchFileName, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0}})
			if err != nil {
				textLogger.Warn("tail file error: ", "err", err, "file", l.WatchFileName)
				return
			}
		}
	}
}

func rangeLogFiles(m map[string]int64) {
	for _, logFile := range logFiles {
		go logFile.Run(m)
	}
}

func searchDir(dirName string, w *sync.WaitGroup) error {
	// 遍历目录
	dir, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range dir {
		fileName := file.Name()
		configName := dirName + "/" + fileName
		if file.IsDir() {
			searchDir(configName, w)
		}
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// 读取文件
		config, err := LoadConfig(configName)
		log.Println(configName, config != nil, err)
		if err != nil {
			continue
		}

		f := &LogConfig{
			ConfigName: fileName,
			Config:     config,
			w:          w,
			Close:      make(chan struct{}, 1),
			sinkWriter: NewSinkWriter(config),
			parser:     NewParser(config),
		}

		logFiles[fileName] = f
	}

	return nil
}
