package code

import (
	"math/rand"
	"time"
)

func Generate(length int) string {
	number := []byte("0123456789")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, number[r.Intn(len(number))])
	}
	return string(result)
}

func Generate6NumberCode() string {
	return Generate(6)
}
