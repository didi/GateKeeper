package zerolog

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

type WFFileWriter struct {
	AccessWriter
}

// NewFileWriter creates and initializes a new AccessWriter.
func NewWFFileWriter(options ...Option) *WFFileWriter {
	w := AccessWriter{}

	for _, opt := range options {
		opt(&w)
	}

	if w.curHour == 0 {
		w.curHour = getNow(time.Now())
	}

	var logFileName, linkName string
	if len(w.FilePrefix) != 0 {
		logFileName = fmt.Sprintf("%s.log.wf.%s", w.FilePrefix, getCurHour(time.Now()))
		if !w.DisableLink {
			linkName = fmt.Sprintf("%s.log.wf", w.FilePrefix)
		}
		fd, err := createFile(logFileName, linkName, &w)
		if err == nil {
			w.FileOut = fd
		}
	}
	wf := WFFileWriter{w}

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

func (w *WFFileWriter) autoClear() {
	if !w.AutoClear {
		return
	}
	now := time.Now()
	expireDate := now.Add(time.Duration(-1*w.ClearHours) * time.Hour)
	d := expireDate.Local()
	beginTime := d.Unix() - int64(d.Minute())*60 - int64(d.Second())

	fs, err := getExpiredFilesByDir(w.FileDir, beginTime, w.FilePrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetFilesByDir[%s] err:%v\n", w.FileDir, err)
	}
	for _, f := range fs {
		os.Remove(f)
	}
}

func (w *WFFileWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	if l < WarnLevel {
		return len(p), nil
	}

	var logFileName, linkName string

	if w.FileOut != nil {
		if w.curHour == 0 {
			w.curHour = getNow(time.Now())
		}

		now := time.Now()      // 当前时间戳
		nowhour := getNow(now) //当前hour 取整时间戳
		if atomic.CompareAndSwapInt64(&w.curHour, compareHourToResult(nowhour > w.curHour, w.curHour), nowhour) {
			cur := getCurHour(now)
			logFileName = fmt.Sprintf("%s.log.wf.%s", w.FilePrefix, cur)
			linkName = fmt.Sprintf("%s.log.wf", w.FilePrefix)
			w.rotate(cur, logFileName, linkName)
			go w.autoClear() //spend time
		}
	} else {
		return len(p), nil
	}

	return w.Write(p)
}
