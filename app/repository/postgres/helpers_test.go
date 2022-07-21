package postgres

import "testing"

func Test_mcontains(t *testing.T) {
	type args struct {
		elems []string
		vv    []string
		algo  containsAlgo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    []string{"2", "3"},
				algo:  exactlyAll,
			},
			want: true,
		},
		{
			name: "empty",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    []string{},
				algo:  exactlyAll,
			},
			want: true,
		},
		{
			name: "nil",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    nil,
				algo:  exactlyAll,
			},
			want: true,
		},
		{
			name: "not_contains",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    []string{"5"},
				algo:  exactlyAll,
			},
			want: false,
		},
		{
			name: "not_contains",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    []string{"1", "2", "3", "4", "5"},
				algo:  exactlyAll,
			},
			want: false,
		},
		{
			name: "at_least_one_contains",
			args: args{
				elems: []string{"1", "2", "3", "4"},
				vv:    []string{"5", "6", "7", "8", "1"},
				algo:  atLeastOne,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mcontains(tt.args.elems, tt.args.vv, tt.args.algo); got != tt.want {
				t.Errorf("mcontains() = %v, want %v", got, tt.want)
			}
		})
	}
}
