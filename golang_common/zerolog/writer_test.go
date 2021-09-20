// +build !binary_log
// +build !windows

package zerolog

import (
	"reflect"
	"testing"
)

func TestMultiSyslogWriter(t *testing.T) {
	sw := &syslogTestWriter{}
	log := New(MultiLevelWriter(SyslogLevelWriter(sw)))
	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")
	log.Log().Msg("nolevel")
	want := []syslogEvent{
		{"Debug", `[DEBUG]||debug` + "\n"},
		{"Info", `[INFO]||info` + "\n"},
		{"Warning", `[WARNING]||warn` + "\n"},
		{"Err", `[ERROR]||error` + "\n"},
		{"Info", `` + "\n"},
	}
	if got := sw.events; !reflect.DeepEqual(got, want) {
		t.Errorf("Invalid syslog message routing: want %v, got %v", want, got)
	}
}
