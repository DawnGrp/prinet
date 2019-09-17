package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHex(t *testing.T) {
	data := []byte{0, 1, 2, 3}
	data = append(data, make([]byte, 10-len(data))...)

	t.Log(data)
}
func String2MD5(plainText string) (md5string string) {

	fmt.Println(plainText)
	m := md5.New()

	md5.Sum([]byte(plainText))

	return hex.EncodeToString(m.Sum(nil))
}
