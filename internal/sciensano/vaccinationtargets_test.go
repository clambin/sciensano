package sciensano_test

import (
	"github.com/clambin/sciensano/internal/sciensano"
	"testing"
)

func TestGetAgeGroupFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{args: args{target: "vaccinations-18-25-first"}, wantOutput: "18-25"},
		{args: args{target: "vaccinations-85+-second"}, wantOutput: "85+"},
		{args: args{target: "vaccinations-85+-second"}, wantOutput: "85+"},
		{args: args{target: "vaccinations--second"}, wantOutput: ""},
		{args: args{target: "invalid"}, wantOutput: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := sciensano.GetAgeGroupFromTarget(tt.args.target); gotOutput != tt.wantOutput {
				t.Errorf("GetAgeGroupFromTarget() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetModeFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name     string
		args     args
		wantMode string
	}{
		{name: "A", args: args{target: "vaccinations-0-16-first"}, wantMode: "A"},
		{name: "B", args: args{target: "vaccinations--second"}, wantMode: "B"},
		{name: "invalid", args: args{target: "foobar"}, wantMode: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMode := sciensano.GetModeFromTarget(tt.args.target); gotMode != tt.wantMode {
				t.Errorf("GetModeFromTarget() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}
