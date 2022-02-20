package utils

import (
	"crypto/sha256"
	"fmt"
)

// https://blog.8bitzen.com/posts/22-08-2019-how-to-hash-a-struct-in-go
func Hash(input interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", input)))

	return fmt.Sprintf("%x", h.Sum(nil))
}