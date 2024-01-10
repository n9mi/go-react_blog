package utils

import "math/rand"

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateRandomString(length int) string {
	res := make([]byte, length)

	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}
