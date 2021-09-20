package zerolog_test

import (
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"os"
	"time"

	"github.com/didi/gatekeeper/golang_common/zerolog"
)

func ExampleNew() {
	log := zerolog.New(os.Stdout)

	log.Info().Msg("hello world")
	// Output: [INFO]||hello world
}

func ExampleLogger_Level() {
	log := zerolog.New(os.Stdout).Level(zerolog.WarnLevel)
	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")

	// Output: [ERROR]||kept message
}

func ExampleLogger_Sample() {
	log := zerolog.New(os.Stdout).Sample(&zerolog.BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	// Output: [INFO]||message 1
	//[INFO]||message 3
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	if l != zerolog.NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook
	var messageHook MessageHook = "The message"

	log := zerolog.New(os.Stdout).Hook(levelNameHook).Hook(messageHook)

	log.Info().Msg("hello world")

	// Output: [INFO]||level_name=info||the_message=hello world||hello world
}

func ExampleLogger_Print() {
	log := zerolog.New(os.Stdout)

	log.Print("hello world")

	// Output: [DEBUG]||hello world
}

func ExampleLogger_Printf() {
	log := zerolog.New(os.Stdout)

	log.Printf("hello %s", "world")

	// Output: [DEBUG]||hello world
}

func ExampleLogger_Debug() {
	log := zerolog.New(os.Stdout)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: [DEBUG]||foo=bar||n=123||hello world
}

func ExampleLogger_Info() {
	log := zerolog.New(os.Stdout)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: [INFO]||foo=bar||n=123||hello world
}

func ExampleLogger_Warn() {
	log := zerolog.New(os.Stdout)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")

	// Output: [WARNING]||foo=bar||a warning message
}

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")

	// Output: [ERROR]||error=some error||error doing something
}

func ExampleLogger_WithLevel() {
	log := zerolog.New(os.Stdout)

	log.WithLevel(zerolog.InfoLevel).
		Msg("hello world")

	// Output: [INFO]||hello world
}

func ExampleLogger_Write() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	// Output: foo=bar||hello world
}

func ExampleLogger_Log() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	// Output: foo=bar||bar=baz
}

// func ExampleEvent_Dict() {
// 	log := zerolog.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Dict("dict", zerolog.Dict().
// 			Str("bar", "baz").
// 			Int("n", 1),
// 		).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
// }

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func (p Price) MarshalZerologObject(e *zerolog.Event) {
	denom := uint64(1)
	for i := 0; i < p.prec; i++ {
		denom *= 10
	}
	result := []byte(p.unit)
	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
	e.Str("price", string(result))
}

type Users []User

func ExampleEvent_Interface() {
	log := zerolog.New(os.Stdout)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")

	// Output: foo=bar||obj={"name":"john"}||hello world
}

func ExampleEvent_Dur() {
	d := time.Duration(10 * time.Second)

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")

	// Output: foo=bar||dur=10000||hello world
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")

	// Output: foo=bar||durs=[10000,20000]||hello world
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger()

	log.Log().Msg("hello world")

	// Output: foo=bar||obj={"name":"john"}||hello world
}

func ExampleContext_Dur() {
	d := time.Duration(10 * time.Second)

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: foo=bar||dur=10000||hello world
}

func ExampleContext_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: foo=bar||durs=[10000,20000]||hello world
}

func ExampleContext_IPAddr() {
	hostIP := net.IP{192, 168, 0, 100}
	log := zerolog.New(os.Stdout).With().
		IPAddr("HostIP", hostIP).
		Logger()

	log.Log().Msg("hello world")

	// Output: HostIP=192.168.0.100||hello world
}

func ExampleContext_IPPrefix() {
	route := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	log := zerolog.New(os.Stdout).With().
		IPPrefix("Route", route).
		Logger()

	log.Log().Msg("hello world")

	// Output: Route=192.168.0.0/24||hello world
}

func ExampleContext_MacAddr() {
	mac := net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	log := zerolog.New(os.Stdout).With().
		MACAddr("hostMAC", mac).
		Logger()

	log.Log().Msg("hello world")

	// Output: hostMAC=00:14:22:01:23:45||hello world
}
