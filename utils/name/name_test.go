package protoparse

import "testing"

func TestCamelCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				s: "",
			},
			want: "",
		},
		{
			args: args{
				s: "name",
			},
			want: "Name",
		},
		{
			args: args{
				s: "test_name",
			},
			want: "TestName",
		},
		{
			args: args{
				s: "NAME",
			},
			want: "NAME",
		},
		{
			args: args{
				s: "gogs_is_testing",
			},
			want: "GogsIsTesting",
		},
		{
			args: args{
				s: "_name",
			},
			want: "XName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelCase(tt.args.s); got != tt.want {
				t.Errorf("CamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
