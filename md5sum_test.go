package md5sum

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const MD5HEX_FOR_EMPTYFILE = "d41d8cd98f00b204e9800998ecf8427e"

func TestChecksumGlob(t *testing.T) {
	var buf bytes.Buffer
	tmpdir, err := ioutil.TempDir(os.TempDir(), "testing-md5sum")
	if err != nil {
		t.Fatal("fail to create directory")
	}
	defer os.RemoveAll(tmpdir)
	tmp1, err := ioutil.TempFile(tmpdir, "testing-md5sum")
	if err != nil {
		t.Fatal("fail to create file")
	}
	tmp2, err := ioutil.TempFile(tmpdir, "testing-md5sum")
	if err != nil {
		t.Fatal("fail to create file")
	}
	if err := ChecksumGlob(fmt.Sprintf("%s/*", tmpdir), &buf); err != nil {
		t.Fatalf("get checksum glob failed: %s\n", err)
	}

	r, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Fatalf("fetch result failed: %s\n", err)
	}
	expected := fmt.Sprintf("%s  %s\n%s  %s\n", MD5HEX_FOR_EMPTYFILE, tmp1.Name(), MD5HEX_FOR_EMPTYFILE, tmp2.Name())
	if string(r) != expected {
		t.Fatalf("not match. expected: %s, actual %s\n", expected, r)
	}
}

func TestChecksumFile(t *testing.T) {
	var buf bytes.Buffer
	tempfile, err := ioutil.TempFile(os.TempDir(), "testing-md5sum")
	if err != nil {
		t.Fatalf("create file failed: %s", err)
	}
	defer os.Remove(tempfile.Name())
	if err := ChecksumFile(tempfile.Name(), &buf); err != nil {
		t.Fatalf("get checksum file failed: %s\n", err)
	}
	r, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Fatalf("fetch result failed: %s\n", err)
	}
	expected := fmt.Sprintf("%s  %s\n", MD5HEX_FOR_EMPTYFILE, tempfile.Name())
	if string(r) != expected {
		t.Fatalf("not match. expected: %s, actual %s\n", expected, r)
	}
}

func TestDecode(t *testing.T) {
	r := strings.NewReader("d41d8cd98f00b204e9800998ecf8427e  hoge\n4af51d184c2507dd9fab8be3766168ac  hoge.md5\n")
	pairs, err := Decode(r)
	if err != nil {
		t.Fatalf("decode failed: %s", err)
	}
	if len(pairs) != 2 {
		t.Fatal("unexpected length of pairs")
	}

	pair1 := &Pair{
		md5sum: "d41d8cd98f00b204e9800998ecf8427e",
		path:   "hoge",
	}
	pair2 := &Pair{
		md5sum: "4af51d184c2507dd9fab8be3766168ac",
		path:   "hoge.md5",
	}
	if !reflect.DeepEqual(pairs[0], pair1) {
		t.Fatalf("field does not match: expected:%v actual: %v", pair1, pairs[0])
	}
	if !reflect.DeepEqual(pairs[1], pair2) {
		t.Fatalf("field does not match: expected:%v actual: %v", pair2, pairs[1])
	}
}

func TestCheck(t *testing.T) {
	var pairs Pairs
	tempfile1, err := ioutil.TempFile(os.TempDir(), "testing-md5sum")
	if err != nil {
		t.Fatalf("create file failed: %s", err)
	}
	defer os.Remove(tempfile1.Name())
	tempfile2, err := ioutil.TempFile(os.TempDir(), "testing-md5sum")
	if err != nil {
		t.Fatalf("create file failed: %s", err)
	}
	defer os.Remove(tempfile2.Name())
	pairs = append(pairs, &Pair{MD5HEX_FOR_EMPTYFILE, tempfile1.Name()})
	pairs = append(pairs, &Pair{MD5HEX_FOR_EMPTYFILE, tempfile2.Name()})
	b, err := Check(pairs)
	if err != nil {
		t.Fatalf("check failed: %s", err)
	}
	if b != true {
		t.Fatal("check does not pass.")
	}
}
