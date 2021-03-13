package apihandler

import (
	"testing"
)

func TestGetAgeGroupFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name         string
		args         args
		wantAgeGroup string
	}{
		{args: args{target: "vaccinations-18-25-first"}, wantAgeGroup: "18-25"},
		{args: args{target: "vaccinations-85+-second"}, wantAgeGroup: "85+"},
		{args: args{target: "vaccinations-85+-second"}, wantAgeGroup: "85+"},
		{args: args{target: "vaccinations--second"}, wantAgeGroup: ""},
		{args: args{target: "invalid"}, wantAgeGroup: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAgeGroup := getAgeGroupFromTarget(tt.args.target); gotAgeGroup != tt.wantAgeGroup {
				t.Errorf("getAgeGroupFromTarget() = %v, want %v", gotAgeGroup, tt.wantAgeGroup)
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
			if gotMode := getModeFromTarget(tt.args.target); gotMode != tt.wantMode {
				t.Errorf("getModeFromTarget() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}

func Test_getRegionFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
	}{
		{name: "Flanders", args: args{target: "vaccinations-Flanders-first"}, wantRegion: "Flanders"},
		{name: "Empty", args: args{target: "vaccinations--second"}, wantRegion: ""},
		{name: "invalid", args: args{target: "foobar"}, wantRegion: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRegion := getRegionFromTarget(tt.args.target); gotRegion != tt.wantRegion {
				t.Errorf("getRegionFromTarget() = %v, want %v", gotRegion, tt.wantRegion)
			}
		})
	}
}
