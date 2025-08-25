package otp

import (
	"crypto/rand"
	"fmt"
	"io"
)

func Generate6() (string, error) {
	var n uint32
	if err := binaryRead(rand.Reader, &n); err != nil {
		return "", err
	}
	// 6 digit
	return fmt.Sprintf("%06d", n%1000000), nil
}

// helper for random
func binaryRead(r io.Reader, out interface{}) error {
	b := make([]byte, 4)
	if _, err := r.Read(b); err != nil {
		return err
	}
	*out.(*uint32) = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return nil
}
