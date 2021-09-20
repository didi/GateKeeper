package zerolog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultTunnelNum = 100000
	defaultWriteSize = 1 << 14
)

type Option func(w *AccessWriter)

// AccessWriter parses the JSON input and writes it in an
// (optionally) filed, human-friendly format to Out.

type AccessWriter struct {
	tunnel        chan []byte
	fileBufWriter *bufio.Writer

	//normal file
	FileOut *os.File

	//old file
	oldFile *os.File

	//auto clear the log
	AutoClear bool

	//mutex for rotate
	mu sync.Mutex

	//当前时间取整小时时间戳
	curHour int64

	//file directory
	FileDir string

	//file prefix name
	FilePrefix string

	//clear hours
	ClearHours int32

	//clear steps
	ClearStep int32

	//is Disable link
	DisableLink bool

	//caller format func
	FormatCallerFunc func(i string) string

	// time format func
	FormatTimestampFunc func() string

	// msg format func
	FormatMessageFunc func(i string) string

	// level format func
	FormatLevelFunc func(i string) string
}

// NewFileWriter creates and initializes a new AccessWriter.
func NewAccessWriter(options ...Option) *AccessWriter {
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

	w.tunnel = make(chan []byte, defaultTunnelNum)
	if w.FileOut != nil {
		w.fileBufWriter = bufio.NewWriterSize(w.FileOut, defaultWriteSize)
	}
	go func() {
		flushTimer := time.NewTimer(time.Second)
		for {
			select {
			case p, ok := <-w.tunnel:
				if !ok {
					return
				}
				w.fileBufWriter.Write(p)
			case <-flushTimer.C:
				if w.fileBufWriter != nil {
					w.fileBufWriter.Flush()
				}
				flushTimer.Reset(time.Second)
			}
		}
	}()
	return &w
}

func createFile(logFileName string, linkName string, w *AccessWriter) (*os.File, error) {
	logFile := filepath.Join(w.FileDir, logFileName)
	if _, err := os.Stat(w.FileDir); os.IsNotExist(err) {
		err = os.Mkdir(w.FileDir, 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create dir %s", err.Error())
			return nil, err
		}
	}

	fd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open file %s \n", err.Error())
		return nil, err
	}

	if !w.DisableLink {
		linkFile := filepath.Join(w.FileDir, linkName)
		os.Remove(linkFile)
		err = os.Symlink(logFile, linkFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create link err: %s \n", err.Error())
			return nil, err
		}
	}

	return fd, nil

}

func SetFileDir(dir string) Option {
	return func(w *AccessWriter) {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get the absolute path %s", err.Error())
			return
		}
		w.FileDir = absPath
	}
}

func SetFilePrefix(prefix string) Option {
	return func(w *AccessWriter) {
		w.FilePrefix = prefix
	}
}

func SetAutoClear(isAutoClear bool) Option {
	return func(w *AccessWriter) {
		w.AutoClear = isAutoClear
	}
}

func SetClearHours(hours int32) Option {
	return func(w *AccessWriter) {
		w.ClearHours = hours
	}
}

func SetClearSteps(steps int32) Option {
	return func(w *AccessWriter) {
		w.ClearStep = steps
	}
}

func SetDisableLink(disableLink bool) Option {
	return func(w *AccessWriter) {
		w.DisableLink = disableLink
	}
}

func getCurHour(cur time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d", cur.Year(), cur.Month(), cur.Day(), cur.Hour())
}

func (w *AccessWriter) rotate(cur string, logFileName string, linkName string) {
	logFile := filepath.Join(w.FileDir, logFileName)
	if w.oldFile != nil {
		w.oldFile.Sync()
		w.oldFile.Close()
	}
	fd, _ := os.Create(logFile)
	w.oldFile = w.FileOut
	w.FileOut = fd
	if w.fileBufWriter != nil {
		w.fileBufWriter.Flush()
	}
	w.fileBufWriter = bufio.NewWriterSize(w.FileOut, defaultWriteSize)
	if !w.DisableLink {
		linkFile := filepath.Join(w.FileDir, linkName)
		os.Remove(linkFile)
		os.Symlink(logFile, linkFile)
	}
}

func (w *AccessWriter) autoClear() {
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
	w.fileBufWriter.Flush()
	for _, f := range fs {
		os.Remove(f)
	}
}

func (w *AccessWriter) WriteLevel(l Level, p []byte) (n int, err error) {

	if l >= WarnLevel {
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

func compareHourToResult(in bool, res int64) int64 {
	if in {
		return res
	}
	return 0
}

func (w *AccessWriter) Write(p []byte) (n int, err error) {
	w.tunnel <- p
	return len(p), nil
}

func (w *AccessWriter) FormatCaller(i string) string {
	if w.FormatCallerFunc != nil {
		return w.FormatCallerFunc(i)
	}
	return defaultFormatCaller(i)
}

func (w *AccessWriter) FormatTimestamp() string {
	if w.FormatTimestampFunc != nil {
		return w.FormatTimestampFunc()
	}
	return defaultFormatTimestamp()
}

func (w *AccessWriter) FormatMessage(i string) string {
	if w.FormatMessageFunc != nil {
		w.FormatMessageFunc(i)
	}
	return defaultFormatMessage(i)
}

func (w *AccessWriter) FormatLevel(i string) string {
	if w.FormatLevelFunc != nil {
		return w.FormatLevelFunc(i)
	}
	return defaultFormatLevel(i)
}

// GetNow 获取当前小时取整后的时间戳
func getNow(cur time.Time) int64 {
	return cur.Unix() - int64(cur.Minute()*60) - int64(cur.Second())
}

// GetNow 获取前一个小时的时间戳
func getHourBefore(cur time.Time) int64 {
	return cur.Unix() - int64(cur.Minute()*60) - int64(cur.Second()) - 3600
}
