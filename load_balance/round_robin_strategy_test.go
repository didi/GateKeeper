package load_balance

import "testing"

func TestRoundRobinStrategy_Add(t *testing.T) {
	type fields struct {
		curIndex int
		rss      []string
	}
	type args struct {
		params []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoundRobinStrategy{
				curIndex: tt.fields.curIndex,
				rss:      tt.fields.rss,
			}
			if err := r.Add(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}