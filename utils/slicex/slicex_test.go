package slicex

import "testing"

func TestInSlice(t *testing.T) {
	type args struct {
		target string
		data   []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				target: "a",
				data:   []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				target: "d",
				data:   []string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "test3",
			args: args{
				target: "d",
				data:   []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InSlice(tt.args.target, tt.args.data); got != tt.want {
				t.Errorf("InSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
