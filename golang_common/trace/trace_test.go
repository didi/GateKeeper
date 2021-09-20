package trace

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetTraceId(t *testing.T) {
	req := httptest.NewRequest("POST", "http://www.test.com", nil)
	tracer := New(req)
	traceid := tracer.GetTraceId(req)
	if len(traceid) <= 0 {
		t.Failed()
	}
	t.Log(traceid)
}

func TestIsTraceSampleEnabled(t *testing.T) {
	req := httptest.NewRequest("POST", "http://www.test.com", nil)
	tracer := New(req)
	hit := tracer.IsTraceSampleEnabled()
	t.Log(hit)
}

func BenchmarkGetTraceId(b *testing.B) {
	req := httptest.NewRequest("POST", "http://www.test.com", nil)
	for i := 0; i < b.N; i++ {
		tracer := New(req)
		tracer.GetTraceId(req)
	}
}

func BenchmarkGetLocalIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getLocalIP()
	}
}

func TestGetCtxTrace(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  *Trace
		want1 bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetCtxTrace(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCtxTrace() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCtxTrace() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewWithMap(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name      string
		args      args
	}{
	// TODO: Add test cases.
		{
			name: "test1",
			args: args {
				m: map[string]string{
					"didi-header-rid":"1234",
					DIDI_HEADER_SPANID:"3435",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTrace := NewWithMap(tt.args.m); gotTrace != nil {
				t.Logf("NewWithMap() = %v", gotTrace)
			}
		})
	}
}
