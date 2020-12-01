package controllers

import (
	"os"
	"testing"
)

func TestGetEnvDefault(t *testing.T) {
	type args struct {
		variable     string
		defaultVal   string
		shouldSetENV bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Non-existing ENV variable",
			args: args{
				variable:   "FOO",
				defaultVal: "BAR",
			},
			want: "BAR",
		},
		{
			name: "Existing ENV variable",
			args: args{
				variable:     "FOO",
				defaultVal:   "BAR",
				shouldSetENV: true,
			},
			want: "FOO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.shouldSetENV {
				_ = os.Setenv(tt.args.variable, tt.args.variable)
			}
			if got := GetEnvDefault(tt.args.variable, tt.args.defaultVal); got != tt.want {
				t.Errorf("GetEnvDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
