package zerolog

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetFilesByDir ...
func getFilesByDir(dir string) ([]string, error) {
	var fileList []string
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dir, f.Name())
		fileList = append(fileList, path)
	}

	return fileList, nil
}

func getExpiredFilesByDir(dir string, beginTime int64, filePrefix string) ([]string, error) {
	var fileList []string
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if !strings.HasPrefix(filename, filePrefix) {
			continue
		}
		idx := strings.LastIndex(filename, ".")
		ts := filename[idx+1:]
		if len(ts) == 10 {
			dateStr := ts[0:8]
			hourStr := ts[8:]
			date, err := time.Parse("20060102", dateStr)
			if err != nil {
				return []string{}, err
			}
			reviseDate := date.Local()
			offset, err := strconv.ParseInt(hourStr, 10, 64)
			if err != nil {
				return []string{}, err
			}
			timeUnix := reviseDate.Unix() - int64(reviseDate.Hour())*3600 + offset*3600
			if timeUnix <= beginTime {
				path := filepath.Join(dir, filename)
				fileList = append(fileList, path)
			}
		}
	}

	return fileList, nil
}
