package rev

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Hash contains the checksum of a file.
type Hash struct {
	hash []byte
}

// NewHash takes a checksum and returns a pointer to a Hash struct.
func NewHash(hash []byte) *Hash {
	return &Hash{
		hash: hash,
	}
}

// NewHashFromFile calculates a hash for the given file path and
// returns a pointer to a Hash struct.
func NewHashFromFile(path string) (*Hash, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer logCloser(f)

	return NewHashFromReader(f)
}

// NewHashFromReader calculates a hash from any io.Reader and
// returns a pointer to a Hash struct.
func NewHashFromReader(reader io.Reader) (*Hash, error) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	return &Hash{hash: h.Sum(nil)}, nil
}

func (c *Hash) String() string {
	return hex.EncodeToString(c.hash)
}

// ShortHash returns the first 10 chars of a hash.
func (c *Hash) ShortHash() string {
	if c.hash != nil {
		return c.String()[:10]
	}
	return ""
}

// HashedPath takes a file path and appends
// a short hash to its file name.
// Returned paths form: /path/to/file-shorthash.ext
func (c *Hash) HashedPath(path string) string {
	hash := c.ShortHash()
	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	if ext == "" {
		return filepath.Join(dir, fmt.Sprintf("%s-%s", file, hash))
	}
	fileNoExt := file[:strings.LastIndex(file, ext)]

	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", fileNoExt, hash, ext))
}
