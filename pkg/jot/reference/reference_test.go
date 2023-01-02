package reference

import (
	"net/url"
	"testing"
)

func TestGetBuildUrlConversion(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Github short path style (no protocol)",
			args: args{
				url: "github.com/rocketblend/official-library/build",
			},
			want: "https://raw.githubusercontent.com/rocketblend/official-library/master/build",
		},
		{
			name: "Github offical path style (no protocol)",
			args: args{
				url: "github.com/rocketblend/official-library/builds/stable/3.2.2",
			},
			want: "https://raw.githubusercontent.com/rocketblend/official-library/master/builds/stable/3.2.2",
		},
		{
			name: "Github (http)",
			args: args{
				url: "http://github.com/rocketblend/official-library/build",
			},
			want: "http://raw.githubusercontent.com/rocketblend/official-library/master/build",
		},
		{
			name: "Github (https)",
			args: args{
				url: "https://github.com/rocketblend/official-library/build",
			},
			want: "https://raw.githubusercontent.com/rocketblend/official-library/master/build",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := convertToUrl(tt.args.url)
			if got.String() != tt.want {
				t.Errorf("NewSourceUrl() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Github",
			args: args{
				url: "http://github.com",
			},
			want: "",
		},
		{
			name: "Gitlab",
			args: args{
				url: "http://gitlab.com",
			},
			want: "invalid host: gitlab.com",
		},
		{
			name: "Google",
			args: args{
				url: "http://google.com",
			},
			want: "invalid host: google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, e := url.Parse(tt.args.url)
			if e != nil {
				t.Errorf("url.Parse() error")
			}

			if err := validateHost(u); err != nil && err.Error() != tt.want {
				t.Errorf("validateHost() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestAddToPathIndex(t *testing.T) {
	type args struct {
		path  string
		index int
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "3 elements",
			args: args{
				path:  "one/two/three",
				index: 1,
				value: "four",
			},
			want: "one/four/two/three",
		},
		{
			name: "4 elements",
			args: args{
				path:  "one/two/three/four",
				index: 3,
				value: "five",
			},
			want: "one/two/three/five/four",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if s := addToPathIndex(tt.args.path, tt.args.index, tt.args.value); s != tt.want {
				t.Errorf("addToPathIndex() = %v, want %v", s, tt.want)
			}
		})
	}
}
