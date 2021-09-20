package zerolog

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: make([]byte, 0, 1<<14),
		}
	},
}

// ErrorMarshalFunc allows customization of global error marshaling
var ErrorMarshalFunc = func(err error) interface{} {
	return err
}

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
type Event struct {
	buf   []byte
	w     LevelWriter
	level Level
	done  func(msg string)
	ch    []Hook // hooks from context
}

func putEvent(e *Event) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(e.buf) > maxSize {
		return
	}
	eventPool.Put(e)
}

func newEvent(w LevelWriter, level Level) *Event {
	e := eventPool.Get().(*Event)
	e.buf = e.buf[:0]
	e.ch = nil
	e.w = w
	e.level = level
	return e
}

func (e *Event) write() (err error) {
	if e == nil {
		return nil
	}
	if e.level != Disabled {
		e.buf = enc.AppendLineBreak(e.buf)
		if e.w != nil {
			_, err = e.w.WriteLevel(e.level, e.buf)
		}
	}
	putEvent(e)
	return
}

// Enabled return false if the *Event is going to be filtered out by
// log level or sampling.
func (e *Event) Enabled() bool {
	return e != nil && e.level != Disabled
}

// Discard disables the event so Msg(f) won't print it.
func (e *Event) Discard() *Event {
	if e == nil {
		return e
	}
	e.level = Disabled
	return nil
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}
	e.msg(msg)
}

// Msgf sends the event with formated msg added as the message field if not empty.
//
// NOTICE: once this methid is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msgf(format string, v ...interface{}) {
	if e == nil {
		return
	}
	e.msg(fmt.Sprintf(format, v...))
}

func (e *Event) msg(msg string) {
	if len(e.ch) > 0 {
		e.ch[0].Run(e, e.level, msg)
		if len(e.ch) > 1 {
			for _, hook := range e.ch[1:] {
				hook.Run(e, e.level, msg)
			}
		}
	}
	if msg != "" && len(e.buf) > 0 {
		if len(e.buf) > 0 {
			e.buf = enc.AppendString(e.buf, "||"+msg)
		} else {
			e.buf = enc.AppendString(e.buf, msg)
		}
	}

	if e.done != nil {
		defer e.done(msg)
	}
	if err := e.write(); err != nil {
		if ErrorHandler != nil {
			ErrorHandler(err)
		} else {
			fmt.Fprintf(os.Stderr, "zerolog: could not write msg:%v event: %v\n", msg, err)
		}
	}
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Event) Fields(fields map[string]interface{}) *Event {
	if e == nil {
		return e
	}
	e.buf = appendFields(e.buf, fields)
	return e
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if e == nil || (key == "" && val == "") {
		return e
	}
	e.buf = enc.AppendString(enc.AppendKey(e.buf, key), val)
	return e
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendStrings(enc.AppendKey(e.buf, key), vals)
	return e
}

// Bytes adds the field key with val as a string to the *Event context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *Event) Bytes(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendBytes(enc.AppendKey(e.buf, key), val)
	return e
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendHex(enc.AppendKey(e.buf, key), val)
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *Event) RawJSON(key string, b []byte) *Event {
	if e == nil {
		return e
	}
	e.buf = appendJSON(enc.AppendKey(e.buf, key), b)
	return e
}

// AnErr adds the field key with serialized err to the *Event context.
// If err is nil, no field is added.
func (e *Event) AnErr(key string, err error) *Event {
	marshaled := ErrorMarshalFunc(err)
	switch m := marshaled.(type) {
	case nil:
		return e
	case error:
		return e.Str(key, m.Error())
	case string:
		return e.Str(key, m)
	default:
		return e.Interface(key, m)
	}
}

// Err adds the field "error" with serialized err to the *Event context.
// If err is nil, no field is added.
// To customize the key name, change zerolog.ErrorFieldName.
func (e *Event) Err(err error) *Event {
	return e.AnErr(ErrorFieldName, err)
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendBool(enc.AppendKey(e.buf, key), b)
	return e
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendBools(enc.AppendKey(e.buf, key), b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInt(enc.AppendKey(e.buf, key), i)
	return e
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInts(enc.AppendKey(e.buf, key), i)
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInt8(enc.AppendKey(e.buf, key), i)
	return e
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInts8(enc.AppendKey(e.buf, key), i)
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInt16(enc.AppendKey(e.buf, key), i)
	return e
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInts16(enc.AppendKey(e.buf, key), i)
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInt32(enc.AppendKey(e.buf, key), i)
	return e
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInts32(enc.AppendKey(e.buf, key), i)
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInt64(enc.AppendKey(e.buf, key), i)
	return e
}

// Tag adds the field key  to the *Event context, wuthout separate flag
func (e *Event) Tag(key string) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendString(e.buf, key)
	return e
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInts64(enc.AppendKey(e.buf, key), i)
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUint(enc.AppendKey(e.buf, key), i)
	return e
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUints(enc.AppendKey(e.buf, key), i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUint8(enc.AppendKey(e.buf, key), i)
	return e
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUints8(enc.AppendKey(e.buf, key), i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUint16(enc.AppendKey(e.buf, key), i)
	return e
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUints16(enc.AppendKey(e.buf, key), i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUint32(enc.AppendKey(e.buf, key), i)
	return e
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUints32(enc.AppendKey(e.buf, key), i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUint64(enc.AppendKey(e.buf, key), i)
	return e
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendUints64(enc.AppendKey(e.buf, key), i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendFloat32(enc.AppendKey(e.buf, key), f)
	return e
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendFloats32(enc.AppendKey(e.buf, key), f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendFloat64(enc.AppendKey(e.buf, key), f)
	return e
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendFloats64(enc.AppendKey(e.buf, key), f)
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Event (or *Context) has one
// already.
func (e *Event) Timestamp() *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendString(e.buf, e.w.FormatTimestamp())
	return e
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendTime(enc.AppendKey(e.buf, key), t, TimeFieldFormat)
	return e
}

// Times adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendTimes(enc.AppendKey(e.buf, key), t, TimeFieldFormat)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendDuration(enc.AppendKey(e.buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendDurations(enc.AppendKey(e.buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	if e == nil {
		return e
	}
	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	e.buf = enc.AppendDuration(enc.AppendKey(e.buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendInterface(enc.AppendKey(e.buf, key), i)
	return e
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
func (e *Event) Caller() *Event {
	return e.caller(CallerSkipFrameCount)
}

//todo default caller
func (e *Event) CallerDefault(caller string) *Event {
	tmp := e.w.FormatCaller(caller)
	e.buf = enc.AppendString(e.buf, tmp)
	return e
}

func (e *Event) caller(skip int) *Event {
	if e == nil {
		return e
	}
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return e
	}
	file = filepath.Base(file)
	name := runtime.FuncForPC(pc).Name()

	if len(name) > maxNameLength {
		name = ".." + name[(len(name)-maxNameLength):]
	}
	var tmp string
	tmp = e.w.FormatCaller(name + "/" + file + ":" + strconv.Itoa(line))
	e.buf = enc.AppendString(e.buf, tmp)

	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendIPAddr(enc.AppendKey(e.buf, key), ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendIPPrefix(enc.AppendKey(e.buf, key), pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendMACAddr(enc.AppendKey(e.buf, key), ha)
	return e
}
