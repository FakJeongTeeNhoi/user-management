package service

import (
	"log"
	"math/rand"
	"strconv"
)

func ParseToUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Println("Failed to parse string to uint: ", err)
	}
	return uint(i)
}

func ParseToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Failed to parse string to int: ", err)
	}
	return i
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
