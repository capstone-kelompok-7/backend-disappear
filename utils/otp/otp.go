package otp

import (
	"math/rand"
	"time"
)

func GenerateRandomOTP(otpLent int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	const n = "0123456789"

	otp := make([]byte, otpLent)
	for i := range otp {
		otp[i] = n[r.Intn(len(n))]
	}

	return string(otp)
}
