package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"testing"
)

func TestVerify(t *testing.T) {
	utf8, e := GbkToUtf8([]byte("fatal error LNK1104: \xce\xde\xb7\xa8\xb4\xf2\xbf\xaa\xce\xc4\xbc\xfe\xa1\xb0D:\\myworkspace\\rust\\rust-sample\\target\\debug\\deps\\simplewebserver.exe\xa1\xb1\r\n"))
	fmt.Printf(string(utf8),e)
}
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}