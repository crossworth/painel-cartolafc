package vk

import (
	"testing"
)

func TestProfileScreenNameOrIDFromURL(t *testing.T) {
	type args struct {
		profileURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{
				profileURL: "http://vk.com/garrocha",
			},
			want:    "garrocha",
			wantErr: false,
		},
		{
			name:    "without protocol",
			args:    args{
				profileURL: "vk.com/brunn0",
			},
			want:    "brunn0",
			wantErr: false,
		},
		{
			name:    "with id",
			args:    args{
				profileURL: "https://vk.com/id271108483",
			},
			want:    "id271108483",
			wantErr: false,
		},
		{
			name:    "mobile link",
			args:    args{
				profileURL: "https://m.vk.com/brunelas",
			},
			want:    "brunelas",
			wantErr: false,
		},
		{
			name:    "with dot",
			args:    args{
				profileURL: "https://vk.com/m.arabe",
			},
			want:    "m.arabe",
			wantErr: false,
		},
		{
			name:    "invalid 1",
			args:    args{
				profileURL: "https://vk.com/invalid-link",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid 2",
			args:    args{
				profileURL: "https://vk.com/invalid_link",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid 3",
			args:    args{
				profileURL: "https://orkut.com/id1010",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProfileScreenNameOrIDFromURL(tt.args.profileURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileScreenNameOrIDFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProfileScreenNameOrIDFromURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
