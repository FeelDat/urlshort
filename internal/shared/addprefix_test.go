package shared

import "testing"

func TestAddPrefix(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid URL without scheme",
			args: args{
				addr: "example.com",
			},
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name: "Valid URL with scheme",
			args: args{
				addr: "https://example.com",
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				addr: "://badurl",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Edge Case - Empty string",
			args: args{
				addr: "",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddPrefix(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddPrefix() got = %v, want %v", got, tt.want)
			}
		})
	}
}
