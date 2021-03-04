package convert

import (
	"crypto/sha1"
	"fmt"
)

var h = sha1.New()

// MakeId takes a byte representation of a resource and returns a stable string ID for it.
func MakeId(s []byte) string {
	h.Reset()
	h.Write(s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
