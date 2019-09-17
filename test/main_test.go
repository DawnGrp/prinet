package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHex(t *testing.T) {
	data := []byte("")
	//data = bytes.TrimSpace(data)
	t.Log(data, string(data))
}
func String2MD5(plainText string) (md5string string) {

	fmt.Println(plainText)
	m := md5.New()

	md5.Sum([]byte(plainText))

	return hex.EncodeToString(m.Sum(nil))
}
