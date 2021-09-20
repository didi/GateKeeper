package zerolog

// nomal file write will not separate for level

import (
	"bufio"
	"fmt"
	"sync/atomic"
	"time"
)

type FileWriter struct {
	AccessWriter
}

// NewFileWriter creates and initializes a new AccessWriter.
func NewFileWriter(options ...Option) *FileWriter {
	w := AccessWriter{}
	for _, opt := range options {
		opt(&w)
	}
	if w.curHour == 0 {
		w.curHour = getNow(time.Now())
	}
	var logFileName, linkName string
	if len(w.FilePrefix) != 0 {
		logFileName = fmt.Sprintf("%s.log.%s", w.FilePrefix, getCurHour(time.Now()))
		if !w.DisableLink {
			linkName = fmt.Sprintf("%s.log", w.FilePrefix)
		}
		fd, err := createFile(logFileName, linkName, &w)
		if err == nil {
			w.FileOut = fd
		}
	}
	wf := FileWriter{w}

	wf.tunnel = make(chan []byte, defaultTunnelNum)
	if w.FileOut != nil {
		wf.fileBufWriter = bufio.NewWriterSize(wf.FileOut, defaultWriteSize)
	}
	go func() {
		flushTimer := time.NewTimer(time.Second)
		for {
			select {
			case p, ok := <-wf.tunnel:
				if !ok {
					return
				}
				wf.fileBufWriter.Write(p)
			case <-flushTimer.C:
				if wf.fileBufWriter != nil {
					wf.fileBufWriter.Flush()
				}
				flushTimer.Reset(time.Second)
			}
		}
	}()
	return &wf
}

//final write
func (w *FileWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	var logFileName, linkName string

	if w.FileOut != nil {
		now := time.Now() // 当前时间戳
		if w.curHour == 0 {
			w.curHour = getNow(time.Now())
		}
		nowhour := getNow(now) //当前hour 取整时间戳
		if atomic.CompareAndSwapInt64(&w.curHour, compareHourToResult(nowhour > w.curHour, w.curHour), nowhour) {
			cur := getCurHour(now)
			logFileName = fmt.Sprintf("%s.log.%s", w.FilePrefix, cur)
			linkName = fmt.Sprintf("%s.log", w.FilePrefix)
			w.rotate(cur, logFileName, linkName)
			go w.autoClear() //spend time
		}
	} else {
		return len(p), nil
	}

	return w.Write(p)
}
