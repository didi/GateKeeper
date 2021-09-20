package zerolog

import (
	"bytes"
	"errors"
	"fmt"

	//	"io/ioutil"
	stdlog "log"
	"time"
)

func ExampleBinaryNew() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Info().Tag(" undefined").Msg("hello world")
	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [INFO] undefined||hello world
}

func ExampleLogger_Level() {
	dst := bytes.Buffer{}
	log := New(&dst).Level(WarnLevel)

	log.Info().Tag(" undefined").Msg("filtered out message")
	log.Error().Tag(" undefined").Msg("kept message")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [ERROR] undefined||kept message
}

func ExampleLogger_Sample() {
	dst := bytes.Buffer{}
	log := New(&dst).Sample(&BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [INFO]||message 1
	//[INFO]||message 3
}

type LevelNameHook1 struct{}

func (h LevelNameHook1) Run(e *Event, l Level, msg string) {
	if l != NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *Event, l Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook1
	var messageHook MessageHook = "The message"

	dst := bytes.Buffer{}
	log := New(&dst).Hook(levelNameHook).Hook(messageHook)

	log.Info().Tag(" undefined").Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [INFO] undefined||level_name=info||the_message=hello world||hello world
}

func ExampleLogger_Print() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Print("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [DEBUG]||hello world
}

func ExampleLogger_Printf() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Printf("hello %s", "world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [DEBUG]||hello world
}

func ExampleLogger_Debug() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Debug().Tag(" undefined").
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [DEBUG] undefined||foo=bar||n=123||hello world
}

func ExampleLogger_Info() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Info().Tag(" undefined").
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [INFO] undefined||foo=bar||n=123||hello world
}

func ExampleLogger_Warn() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Warn().Tag(" undefined").
		Str("foo", "bar").
		Msg("a warning message")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// [WARNING] undefined||foo=bar||a warning message
}

func ExampleLogger_Error() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Error().Tag(" undefined").
		Err(errors.New("some error")).
		Msg("error doing something")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [ERROR] undefined||error=some error||error doing something
}

func ExampleLogger_WithLevel() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.WithLevel(InfoLevel).Tag(" undefined").
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: [INFO] undefined||hello world
}

func ExampleLogger_Write() {
	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: foo=bar||hello world
}

func ExampleLogger_Log() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().Tag(" undefined").
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: undefined||foo=bar||bar=baz
}

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Users []User

func ExampleEvent_Interface() {
	dst := bytes.Buffer{}
	log := New(&dst)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().Tag(" undefined").
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: undefined||foo=bar||obj={"name":"john"}||hello world
}

func ExampleEvent_Dur() {
	d := time.Duration(10 * time.Second)

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().Tag(" undefined").
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: undefined||foo=bar||dur=10000||hello world
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().Tag(" undefined").
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: undefined||foo=bar||durs=[10000,20000]||hello world
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: foo=bar||obj={"name":"john"}||hello world
}

func ExampleContext_Dur() {
	d := time.Duration(10 * time.Second)

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: foo=bar||dur=10000||hello world
}

func ExampleContext_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: foo=bar||durs=[10000,20000]||hello world
}
