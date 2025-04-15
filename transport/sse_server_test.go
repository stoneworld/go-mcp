package transport

import (
	"net/url"
	"testing"
)

func Test_joinPath(t *testing.T) {
	type args struct {
		u    *url.URL
		elem []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				u: func() *url.URL {
					uri, err := url.Parse("https://google.com/api/v1")
					if err != nil {
						panic(err)
					}
					return uri
				}(),
				elem: []string{"/test"},
			},
			want: "https://google.com/api/v1/test",
		},
		{
			name: "2",
			args: args{
				u: func() *url.URL {
					uri, err := url.Parse("/api/v1")
					if err != nil {
						panic(err)
					}
					return uri
				}(),
				elem: []string{"/test"},
			},
			want: "/api/v1/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinPath(tt.args.u, tt.args.elem...)
			if got := tt.args.u.String(); got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
