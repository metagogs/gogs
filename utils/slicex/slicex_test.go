package slicex

import (
	"reflect"
	"testing"
)

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

	type argsInt struct {
		target int
		data   []int
	}
	testsInt := []struct {
		name string
		args argsInt
		want bool
	}{
		{
			name: "test1",
			args: argsInt{
				target: 1,
				data:   []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "test2",
			args: argsInt{
				target: 4,
				data:   []int{1, 2, 3},
			},
			want: false,
		},
		{
			name: "test3",
			args: argsInt{
				target: 1,
				data:   []int{},
			},
			want: false,
		},
	}
	for _, tt := range testsInt {
		t.Run(tt.name, func(t *testing.T) {
			if got := InSlice(tt.args.target, tt.args.data); got != tt.want {
				t.Errorf("InSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveSliceItem(t *testing.T) {
	type args struct {
		list []string
		item string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				list: []string{"a", "b", "c"},
				item: "b",
			},
			want: []string{"a", "c"},
		},
		{
			name: "test2",
			args: args{
				list: []string{"a", "b", "c"},
				item: "d",
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "test3",
			args: args{
				list: []string{},
				item: "d",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSliceItem(tt.args.list, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveSliceItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
