package service

import (
	"log"
	"strconv"
)

func ParseToUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Println("Failed to parse string to uint: ", err)
	}
	return uint(i)
}
