package rev

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// File contains all information for an asset.
type File struct {
	Name         string      // logo.png
	Ext          string      // .png
	Rel          string      // img/logo.png
	Path         string      // /home/dev/project/public/img/logo.png
	URL          string      // https://www.randomz.com/img/logo.png
	Size         int64       // 512 bytes
	Hash         string      // 2a06ac32a7
	HashedRel    string      // img/logo-2a06ac32a7.png
	HashedPath   string      // /home/dev/project/public/img/logo-2a06ac32a7.png
	HashedURL    string      // https://www.randomz.com/img/logo-2a06ac32a7.png
	ReplacedRefs int         // counts amount of replaced references
	baseurl      string      // https://www.randomz.com/
	basepath     string      // /home/dev/project/public
	filemode     os.FileMode // perm
	hash         *Hash       // (initial) hash of file
}

// NewFile returns a File.
// Returns error on invalid paths.
func NewFile(file, basepath, baseurl string, size int64, mode os.FileMode) (*File, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	base, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return nil, err
	}

	f := &File{
		Name:     name,
		Ext:      ext,
		Rel:      rel,
		Path:     path,
		URL:      baseurl + rel,
		Size:     size,
		baseurl:  baseurl,
		basepath: base,
		filemode: mode,
	}
	return f, nil
}

// ReplaceReferences reads the content of the given file and
// replaces any occurences of referenced files. Changed file contents
// are written in place to the same file.
// If no references are found in the file's content, no chages occur.
func (f *File) ReplaceReferences(refs Files) error {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	counter := 0
	for _, ref := range refs {
		counter += bytes.Count(buf, []byte(ref.URL))
		buf = bytes.Replace(buf, []byte(ref.URL), []byte(ref.HashedURL), -1)
	}
	f.ReplacedRefs = counter
	if counter > 0 {
		return ioutil.WriteFile(f.Path, buf, f.filemode)
	}
	return nil
}

// Checksum creates a hash of the given file.
func (f *File) Checksum() error {
	hash, err := NewHashFromFile(f.Path)
	if err != nil {
		return err
	}
	f.hash = hash
	f.Hash = hash.ShortHash()
	f.HashedRel = hash.HashedPath(f.Rel)
	f.HashedPath = hash.HashedPath(f.Path)
	f.HashedURL = f.baseurl + f.HashedRel

	return nil
}

// Revision creates a hash of the given file,
// appends the hash to the file name and renames
// the file.
func (f *File) Revision() error {
	if err := f.Checksum(); err != nil {
		return err
	}
	oldpath := f.Path
	newpath := f.HashedPath

	return os.Rename(oldpath, newpath)
}

// Files contains multiple files.
type Files []File

// SortByPath sorts files by their path in ascending order.
func (f Files) SortByPath() {
	sort.Slice(f, func(i, j int) bool { return f[i].Path < f[j].Path })
}
