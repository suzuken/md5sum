package md5sum

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ErrPathNotFound = errors.New("file or directory not found")

// Pair is each line of md5sum
type Pair struct {
	md5sum string // hex string of checksum
	path   string // the name of original file
}

// Pairs are multiple
type Pairs []*Pair

func ChecksumGlob(pattern string, w io.Writer) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return ErrPathNotFound
	}
	for _, m := range matches {
		if err := ChecksumFile(m, w); err != nil {
			return err
		}
	}
	return nil
}

// ChechsumFile writes checksum and its file path into writer.
func ChecksumFile(path string, w io.Writer) error {
	pair, err := Calc(path)
	if err != nil {
		return err
	}
	if _, err := WritePair(w, pair); err != nil {
		return err
	}
	return nil
}

// Calc creates pair of md5 checksum and its file path
func Calc(path string) (*Pair, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Pair{
		md5sum: fmt.Sprintf("%x", md5.Sum(b)),
		path:   path,
	}, nil
}

// Encode encodes checksum into writer
func WritePair(w io.Writer, pair *Pair) (n int, err error) {
	return fmt.Fprintf(w, "%s  %s\n", pair.md5sum, pair.path)
}

// Decode decodes chechsum stream into field
func Decode(r io.Reader) (Pairs, error) {
	var pairs Pairs
	s := bufio.NewScanner(r)
	// TODO: should do atomically
	for s.Scan() {
		ss := strings.Split(s.Text(), "  ") // separated by 2 spaces
		if len(ss) != 2 {
			continue
		}
		pair := &Pair{
			md5sum: ss[0],
			path:   ss[1],
		}
		pairs = append(pairs, pair)
	}
	if err := s.Err(); err != nil {
		return Pairs{}, err
	}
	return pairs, nil
}

// Check verify if given file is match
func CheckGlob(pattern string, w io.Writer) (bool, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false, err
	}
	if len(matches) == 0 {
		return false, ErrPathNotFound
	}
	for _, m := range matches {
		if b, err := CheckFile(m); err != nil {
			return false, err
		} else if b == false {
			return false, err
		}
	}
	return true, nil
}

// ReadChecksumFile read checksum file and return pairs
func ReadChecksumFile(path string) (Pairs, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Decode(f)
}

// CheckFile receives md5sum file, and check each line
func CheckFile(checksumFilePath string) (bool, error) {
	pairs, err := ReadChecksumFile(checksumFilePath)
	if err != nil {
		return false, err
	}
	return Check(pairs)
}

// Check checking files using md5 checksum file. For example,
//
// 4af51d184c2507dd9fab8be3766168ac  hoge
//
// checks file `hoge`, and calculate md5 checksum thereof.
// If it equals 4af51d184c2507dd9fab8be3766168ac, return true.
func Check(pairs Pairs) (bool, error) {
	for _, p := range pairs {
		// read path, calculate md5, and test it
		pair, err := Calc(p.path)
		if err != nil {
			return false, err
		}
		if pair.md5sum != p.md5sum {
			return false, nil
		}
	}
	return true, nil
}
