package domain

import (
	"fmt"
	"math/rand"
)

func RandomFromSlice[T interface{}](slice []T) T {
	switch len(slice) {
	case 0:
		panic(fmt.Sprintf("got empty slice %T", slice))
	case 1:
		return slice[0]
	default:
		return slice[rand.Intn(len(slice)-1)]
	}
}
