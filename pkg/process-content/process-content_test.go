package process_content

import (
	"testing"
)

func Test_hasContentChanged(t *testing.T) {
	type args struct {
		contentA string
		contentB string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "content has changed",
			args: args{
				contentA: "1",
				contentB: "2",
			},
			want: true,
		},
		{
			name: "content has not changed",
			args: args{
				contentA: "1",
				contentB: "1",
			},
			want: false,
		},
		{
			name: "content has not changed",
			args: args{
				contentA: "",
				contentB: "",
			},
			want: false,
		},
		{
			name: "content has changed",
			args: args{
				contentA: "xxxx",
				contentB: "xxx",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasContentChanged(tt.args.contentA, tt.args.contentB); got != tt.want {
				t.Errorf("hasContentChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}
