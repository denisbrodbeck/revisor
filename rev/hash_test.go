package rev

import "testing"

func TestHash_HashedPath(t *testing.T) {
	type fields struct {
		hash []byte
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "normal",
			fields: fields{hash: []byte("1234567890")},
			args:   args{path: "/some/random/file.png"},
			want:   "/some/random/file-1234567890.png",
		},
		{
			name:   "no-ext",
			fields: fields{hash: []byte("1234567890")},
			args:   args{path: "/some/random/file"},
			want:   "/some/random/file-1234567890",
		},
		{
			name:   "no-file",
			fields: fields{hash: []byte("1234567890")},
			args:   args{path: "/some/random/"},
			want:   "/some/random/-1234567890", // this is an invalid path upstream
		},
		{
			name:   "no-hash",
			fields: fields{hash: []byte("")},
			args:   args{path: "/some/random/file.png"},
			want:   "/some/random/file-.png", // this is an invalid path upstream
		},
		{
			name:   "no-path",
			fields: fields{hash: []byte("1234567890")},
			args:   args{path: "file.png"},
			want:   "file-1234567890.png",
		},
		{
			name:   "empty",
			fields: fields{hash: []byte("")},
			args:   args{path: ""},
			want:   "-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Hash{
				hash: tt.fields.hash,
			}
			if got := c.HashedPath(tt.args.path); got != tt.want {
				t.Errorf("Hash.HashedPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
