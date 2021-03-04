package convert

import (
	"crypto/sha1"
	"fmt"
)

var h = sha1.New()

func MakeId(s []byte) string {
	h.Reset()
	h.Write(s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
