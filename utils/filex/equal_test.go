package filex

import "testing"

func TestIsFileEqual(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "euqal file",
			args: args{
				a: "testdata/a_go.mod",
				b: "testdata/b_go.mod",
			},
			want: true,
		},
		{
			name: "not euqal file",
			args: args{
				a: "testdata/a_go.mod",
				b: "testdata/c_go.mod",
			},
			want: false,
		},
		{
			name: "not euqal file",
			args: args{
				a: "testdata/e_go.mod",
				b: "testdata/c_go.mod",
			},
			want: false,
		},
		{
			name: "not euqal file",
			args: args{
				a: "testdata/a_go.mod",
				b: "testdata/f_go.mod",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("IsFileEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
