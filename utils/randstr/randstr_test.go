package randstr

import "testing"

func TestRandStr(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "rand",
			args: args{
				n: 17,
			},
			want: 17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStr(tt.args.n); len(got) != tt.want {
				t.Errorf("RandStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
